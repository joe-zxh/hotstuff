package consensus

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/relab/hotstuff/config"
	"github.com/relab/hotstuff/data"
	"github.com/relab/hotstuff/internal/logging"
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// EventType is the type of notification sent to pacemaker
type EventType uint8

// These are the types of events that can be sent to pacemaker
const (
	HQCUpdate EventType = iota
)

// HotStuffCore is the safety core of the HotStuffCore protocol
type HotStuffCore struct {
	Mut sync.Mutex

	// Contains the commands that are waiting to be proposed
	cmdCache *data.CommandSet
	Config   *config.ReplicaConfig
	Blocks   data.BlockStorage
	SigCache *data.SignatureCache

	// protocol data
	vHeight    int
	genesis    *data.Block
	bLock      *data.Block // 已经precommit的block
	bExec      *data.Block // 已经commit的block
	BLeaf      *data.Block // 叶节点对应的block
	qcHigh     *data.QuorumCert
	PendingQCs map[data.BlockHash]*data.QuorumCert

	waitProposal *sync.Cond // 这里的条件变量可以借鉴一下

	pendingUpdates chan *data.Block

	// stops any goroutines started by HotStuff
	cancel context.CancelFunc

	exec chan []data.Command
}

func (hs *HotStuffCore) AddCommand(command data.Command) {
	hs.cmdCache.Add(command)
}

// GetHeight returns the height of the tree
func (hs *HotStuffCore) GetHeight() int {
	return hs.BLeaf.Height
}

// GetVotedHeight returns the height that was last voted at
func (hs *HotStuffCore) GetVotedHeight() int {
	return hs.vHeight
}

// GetLeaf returns the current leaf node of the tree
func (hs *HotStuffCore) GetLeaf() *data.Block {
	hs.Mut.Lock()
	defer hs.Mut.Unlock()
	return hs.BLeaf
}

// SetLeaf sets the leaf node of the tree
func (hs *HotStuffCore) SetLeaf(block *data.Block) {
	hs.Mut.Lock()
	defer hs.Mut.Unlock()
	hs.BLeaf = block
}

// GetQCHigh returns the highest valid Quorum Certificate known to the hotstuff instance.
func (hs *HotStuffCore) GetQCHigh() *data.QuorumCert {
	hs.Mut.Lock()
	defer hs.Mut.Unlock()
	return hs.qcHigh
}

func (hs *HotStuffCore) GetExec() chan []data.Command {
	return hs.exec
}

// New creates a new Hotstuff instance
func New(conf *config.ReplicaConfig) *HotStuffCore {
	logger.SetPrefix(fmt.Sprintf("hs(id %d): ", conf.ID))
	genesis := &data.Block{
		Committed: true,
	}
	qcForGenesis := data.CreateQuorumCert(genesis)
	blocks := data.NewMapStorage()
	blocks.Put(genesis)

	ctx, cancel := context.WithCancel(context.Background())

	hs := &HotStuffCore{
		Config:         conf,
		genesis:        genesis,
		bLock:          genesis,
		bExec:          genesis,
		BLeaf:          genesis,
		qcHigh:         qcForGenesis,
		Blocks:         blocks,
		PendingQCs:     make(map[data.BlockHash]*data.QuorumCert),
		cancel:         cancel,
		SigCache:       data.NewSignatureCache(conf),
		cmdCache:       data.NewCommandSet(),
		pendingUpdates: make(chan *data.Block, 1),
		exec:           make(chan []data.Command, 1),
	}

	hs.waitProposal = sync.NewCond(&hs.Mut)

	go hs.updateAsync(ctx)

	return hs
}

// expectBlock waits for a block with the given Hash
// hs.mut must be locked when calling this function
func (hs *HotStuffCore) ExpectBlock(hash data.BlockHash) (*data.Block, bool) {
	count := 0
	for {
		if block, ok := hs.Blocks.Get(hash); ok {
			return block, true
		} else {
			count++
			hs.waitProposal.Wait()
		}
		if count == 10 {
			log.Printf("warning: waiting for block(%.8s) over 10 times...\n", hash)
			count = 0
		}
	}
}

// UpdateQCHigh updates the qc held by the paceMaker, to the newest qc.
func (hs *HotStuffCore) UpdateQCHigh(qc *data.QuorumCert) bool {
	if !hs.SigCache.VerifyQuorumCert(qc) {
		log.Println("QC not verified!:", qc)
		return false
	}

	logger.Println("UpdateQCHigh")

	newQCHighBlock, ok := hs.ExpectBlock(qc.BlockHash)
	if !ok {
		log.Println("Could not find block of new QC!")
		return false
	}

	oldQCHighBlock, ok := hs.Blocks.BlockOf(hs.qcHigh)
	if !ok {
		panic(fmt.Errorf("Block from the old qcHigh missing from storage"))
	}

	if newQCHighBlock.Height > oldQCHighBlock.Height {
		hs.qcHigh = qc
		hs.BLeaf = newQCHighBlock
		return true
	}

	logger.Println("UpdateQCHigh Failed")
	return false
}

// OnReceiveProposal handles a replica's response to the Proposal from the leader
func (hs *HotStuffCore) OnReceiveProposal(block *data.Block) (*data.PartialCert, error) {
	logger.Println("OnReceiveProposal:", block)
	hs.Blocks.Put(block)
	hs.waitProposal.Broadcast()

	hs.Mut.Lock()
	qcBlock, nExists := hs.ExpectBlock(block.Justify.BlockHash)

	if block.Height <= hs.vHeight {
		hs.Mut.Unlock()
		log.Printf("OnReceiveProposal: Block height(%d) less than vHeight(%d)\n", block.Height, hs.vHeight)
		return nil, fmt.Errorf("Block was not accepted")
	}

	safe := false
	if nExists && qcBlock.Height > hs.bLock.Height {
		safe = true
	} else {
		logger.Println("OnReceiveProposal: liveness condition failed")
		// check if block extends bLock
		b := block
		ok := true
		for ok && b.Height > hs.bLock.Height+1 {
			b, ok = hs.Blocks.Get(b.ParentHash)
		}
		if ok && b.ParentHash == hs.bLock.Hash() {
			safe = true
		} else {
			log.Println("OnReceiveProposal: safety condition failed")
		}
	}

	if !safe {
		hs.Mut.Unlock()
		log.Println("OnReceiveProposal: Block not safe")
		return nil, fmt.Errorf("Block was not accepted")
	}

	logger.Println("OnReceiveProposal: Accepted block")
	hs.vHeight = block.Height
	hs.cmdCache.MarkProposed(block.Commands...)
	hs.Mut.Unlock()

	hs.pendingUpdates <- block

	pc, err := hs.SigCache.CreatePartialCert(hs.Config.ID, hs.Config.PrivateKey, block)
	if err != nil {
		panic(err)
		return nil, err
	}
	return pc, nil
}

func (hs *HotStuffCore) updateAsync(ctx context.Context) {
	for {
		select {
		case n := <-hs.pendingUpdates:
			hs.update(n)
		case <-ctx.Done():
			return
		}
	}
}

func (hs *HotStuffCore) update(block *data.Block) {
	// block1 = b'', block2 = b', block3 = b
	block1, ok := hs.Blocks.BlockOf(block.Justify)
	if !ok || block1.Committed {
		return
	}

	hs.Mut.Lock()
	defer hs.Mut.Unlock()

	logger.Println("PRE COMMIT:", block1)
	// PRE-COMMIT on block1
	hs.UpdateQCHigh(block.Justify)

	block2, ok := hs.Blocks.BlockOf(block1.Justify)
	if !ok || block2.Committed {
		return
	}

	if block2.Height > hs.bLock.Height {
		hs.bLock = block2 // COMMIT on block2
		logger.Println("COMMIT:", block2)
	}

	block3, ok := hs.Blocks.BlockOf(block2.Justify)
	if !ok || block3.Committed {
		return
	}

	if block1.ParentHash == block2.Hash() && block2.ParentHash == block3.Hash() {
		logger.Println("DECIDE", block3)
		hs.commit(block3)
		hs.bExec = block3 // DECIDE on block3
	}

	// Free up space by deleting old data
	hs.Blocks.GarbageCollectBlocks(hs.GetVotedHeight())
	hs.cmdCache.TrimToLen(hs.Config.BatchSize * 5)
	hs.SigCache.EvictOld(hs.Config.QuorumSize * 5)
}

func (hs *HotStuffCore) commit(block *data.Block) {
	// only called from within update. Thus covered by its mutex lock.
	if hs.bExec.Height < block.Height {
		if parent, ok := hs.Blocks.ParentOf(block); ok {
			hs.commit(parent) // todo: 递归改循环
		}
		block.Committed = true
		logger.Println("EXEC", block)
		hs.exec <- block.Commands
	}
}

// CreateProposal creates a new proposal
func (hs *HotStuffCore) CreateProposal() *data.Block {
	batch := hs.cmdCache.GetFirst(hs.Config.BatchSize)
	hs.Mut.Lock()
	b := CreateLeaf(hs.BLeaf, batch, hs.qcHigh, hs.BLeaf.Height+1)
	hs.Mut.Unlock()
	b.Proposer = hs.Config.ID
	hs.Blocks.Put(b)
	return b
}

// Close frees resources held by HotStuff and closes backend connections
func (hs *HotStuffCore) Close() {
	hs.cancel()
}

// CreateLeaf returns a new block that extends the parent.
func CreateLeaf(parent *data.Block, cmds []data.Command, qc *data.QuorumCert, height int) *data.Block {
	return &data.Block{
		ParentHash: parent.Hash(),
		Commands:   cmds,
		Justify:    qc,
		Height:     height,
	}
}
