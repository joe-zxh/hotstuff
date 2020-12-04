package hotstuff

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"net"
	"sync"
	"time"

	"github.com/joe-zxh/hotstuff/config"
	"github.com/joe-zxh/hotstuff/consensus"
	"github.com/joe-zxh/hotstuff/data"
	"github.com/joe-zxh/hotstuff/internal/logging"
	"github.com/joe-zxh/hotstuff/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// Pacemaker is a mechanism that provides synchronization
type Pacemaker interface {
	GetLeader(view int) config.ReplicaID
}

// HotStuff is a thing
type HotStuff struct {
	*consensus.HotStuffCore
	proto.UnimplementedHotstuffServer

	tls bool

	pacemaker Pacemaker

	nodes map[config.ReplicaID]*proto.HotstuffClient
	conns map[config.ReplicaID]*grpc.ClientConn

	server *hotstuffServer

	closeOnce sync.Once

	qcTimeout      time.Duration
	connectTimeout time.Duration
}

//New creates a new GorumsHotStuff backend object.
func New(conf *config.ReplicaConfig, pacemaker Pacemaker, tls bool, connectTimeout, qcTimeout time.Duration) *HotStuff {
	hs := &HotStuff{
		pacemaker:      pacemaker,
		HotStuffCore:   consensus.New(conf),
		nodes:          make(map[config.ReplicaID]*proto.HotstuffClient),
		conns:          make(map[config.ReplicaID]*grpc.ClientConn),
		connectTimeout: connectTimeout,
		qcTimeout:      qcTimeout,
	}
	return hs
}

//Start starts the server and client
func (hs *HotStuff) Start() error {
	addr := hs.Config.Replicas[hs.Config.ID].Address
	err := hs.startServer(addr)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Server: %w", err)
	}
	err = hs.startClient(hs.connectTimeout)
	if err != nil {
		return fmt.Errorf("Failed to start GRPC Clients: %w", err)
	}
	return nil
}

// 作为rpc的client端，调用其他hsserver的rpc。
func (hs *HotStuff) startClient(connectTimeout time.Duration) error {

	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithReturnConnectionError(),
	}

	if hs.tls {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(hs.Config.CertPool, "")))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}

	for rid, replica := range hs.Config.Replicas {
		if replica.ID != hs.Config.ID {
			conn, err := grpc.Dial(replica.Address, grpcOpts...)
			if err != nil {
				log.Fatalf("connect error: %v", err)
				conn.Close()
			} else {
				hs.conns[rid] = conn
				c := proto.NewHotstuffClient(conn)
				hs.nodes[rid] = &c
			}
		}
	}

	return nil
}

// startServer runs a new instance of hotstuffServer
func (hs *HotStuff) startServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Failed to listen to port %s: %w", port, err)
	}

	grpcServerOpts := []grpc.ServerOption{}

	if hs.tls {
		grpcServerOpts = append(grpcServerOpts, grpc.Creds(credentials.NewServerTLSFromCert(hs.Config.Cert)))
	}

	hs.server = newHotStuffServer(hs)

	s := grpc.NewServer(grpcServerOpts...)
	proto.RegisterHotstuffServer(s, hs.server)

	go s.Serve(lis)
	return nil
}

// Close closes all connections made by the HotStuff instance
func (hs *HotStuff) Close() {
	hs.closeOnce.Do(func() {
		hs.HotStuffCore.Close()
		for _, conn := range hs.conns { // close clients connections
			conn.Close()
		}
	})
}

// Propose broadcasts a new proposal to all replicas
func (hs *HotStuff) Propose() {
	hs.Mut.Lock()

	if hs.IsVCExp {
		if hs.ViewChangeCount == 0 {
			hs.VCStart = time.Now()
		}

		if hs.IsVCExp && hs.ViewChangeCount >= hs.ViewChangeMyTotal {
			log.Printf("view change average time: %vms\n", float64(time.Since(hs.VCStart).Milliseconds())/float64(hs.server.ClusterSize*hs.ViewChangeMyTotal))

			<-hs.ViewChangeChan //不再执行了停止了...
		}
		hs.ViewChangeCount++
	}

	proposal := hs.CreateProposal()
	logger.Printf("Propose (%d commands): %s\n", len(proposal.Commands), proposal)
	hs.Mut.Unlock()
	protobuf := proto.BlockToProto(proposal)

	go func() {
		logger.Printf("[B/Propose]: proposer id: %d\n", hs.Config.ID)

		for rid, client := range hs.nodes {
			if rid != hs.Config.ID {
				go func(id config.ReplicaID, cli *proto.HotstuffClient) {
					_, err := (*cli).Propose(context.TODO(), protobuf)
					if err != nil {
						panic(err)
					}
				}(rid, client)
			}
		}
	}()

	hs.handlePropose(proposal)
}

// 这个hotstuffServer是面向 集群内部的，也就是面向hotstuffServer。
type hotstuffServer struct {
	*HotStuff

	mut     sync.RWMutex
	clients map[context.Context]config.ReplicaID
}

func newHotStuffServer(hs *HotStuff) *hotstuffServer {
	hsSrv := &hotstuffServer{
		HotStuff: hs,
		clients:  make(map[context.Context]config.ReplicaID),
	}
	return hsSrv
}

func (hs *HotStuff) handlePropose(block *data.Block) {
	p, err := hs.OnReceiveProposal(block)
	if err != nil {
		logger.Println("OnReceiveProposal returned with error:", err)
		return
	}
	leaderID := hs.pacemaker.GetLeader(block.Height)
	if hs.Config.ID == leaderID {
		hs.handleVote(p)
	} else if leader, ok := hs.nodes[leaderID]; ok {
		(*leader).Vote(context.TODO(), proto.PartialCertToProto(p))
	}
}

// Propose handles a replica's response to the Propose QC from the leader
func (hs *hotstuffServer) Propose(ctx context.Context, protoB *proto.Block) (*empty.Empty, error) {
	block := protoB.FromProto()
	hs.handlePropose(block)
	return &empty.Empty{}, nil
}

// handleVote handles an incoming vote from a replica
func (hs *HotStuff) handleVote(cert *data.PartialCert) {
	if !hs.SigCache.VerifySignature(cert.Sig, cert.BlockHash) {
		log.Println("handleVote: signature not verified!")
		return
	}

	logger.Printf("handleVote: %.8s\n", cert.BlockHash)

	hs.Mut.Lock()
	b := hs.GetBlock(cert.BlockHash)

	qc, ok := hs.PendingQCs[cert.BlockHash]
	if !ok {
		if b.Height <= hs.BLeaf.Height {
			// too old, don't care。已经 这个block已经有qc了
			hs.Mut.Unlock()
			return
		}
		// need to check again in case a qc was created while we waited for the block
		qc, ok = hs.PendingQCs[cert.BlockHash]
		if !ok {
			qc = data.CreateQuorumCert(b)
			hs.PendingQCs[cert.BlockHash] = qc
		}
	}

	err := qc.AddPartial(cert)
	if err != nil {
		panic(err)
		logger.Println("handleVote: could not add partial signature to QC:", err)
	}

	if len(qc.Sigs) >= hs.Config.QuorumSize {
		delete(hs.PendingQCs, cert.BlockHash)
		logger.Printf("handleVote: Created QC: %.8s\n", qc.BlockHash)
		if hs.UpdateQCHigh(qc) {
			go hs.Propose()
		}
	}

	hs.DeletePendings()
	hs.Mut.Unlock()
}

func (hs *hotstuffServer) Vote(_ context.Context, pC *proto.PartialCert) (*empty.Empty, error) {
	go hs.handleVote(pC.FromProto())
	return &empty.Empty{}, nil
}
