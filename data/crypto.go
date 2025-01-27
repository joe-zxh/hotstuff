package data

import (
	"bytes"
	"container/list"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/joe-zxh/hotstuff/config"
	"github.com/joe-zxh/hotstuff/internal/logging"
)

var logger *log.Logger

func init() {
	logger = logging.GetLogger()
}

// SignatureCache keeps a cache of verified signatures in order to speed up verification
type SignatureCache struct {
	conf               *config.ReplicaConfig
	verifiedSignatures map[string]bool
	cache              list.List
	mut                sync.Mutex
}

// NewSignatureCache returns a new instance of SignatureVerifier
func NewSignatureCache(conf *config.ReplicaConfig) *SignatureCache {
	return &SignatureCache{
		conf:               conf,
		verifiedSignatures: make(map[string]bool),
	}
}

// CreatePartialCert creates a partial cert from a block.
func (s *SignatureCache) CreatePartialCert(id config.ReplicaID, privKey *ecdsa.PrivateKey, block *Block) (*PartialCert, error) {
	hash := block.Hash()
	R, S, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}
	sig := PartialSig{id, R, S}
	k := string(sig.ToBytes())
	s.mut.Lock()
	s.verifiedSignatures[k] = true
	s.cache.PushBack(k)
	s.mut.Unlock()
	return &PartialCert{sig, hash}, nil
}

// VerifySignature verifies a partial signature
func (s *SignatureCache) VerifySignature(sig PartialSig, hash BlockHash) bool {
	k := string(sig.ToBytes())

	s.mut.Lock()
	if valid, ok := s.verifiedSignatures[k]; ok {
		s.mut.Unlock()
		return valid
	}
	s.mut.Unlock()

	info, ok := s.conf.Replicas[sig.ID]
	if !ok {
		return false
	}
	valid := ecdsa.Verify(info.PubKey, hash[:], sig.R, sig.S)

	s.mut.Lock()
	s.cache.PushBack(k)
	s.verifiedSignatures[k] = valid
	s.mut.Unlock()

	return valid
}

// VerifyQuorumCert verifies a quorum certificate
func (s *SignatureCache) VerifyQuorumCert(qc *QuorumCert) bool {
	if len(qc.Sigs) < s.conf.QuorumSize {
		return false
	}
	//****
	for _, psig := range qc.Sigs { // 因为需要深度拷贝，所以用range的方式来做，只检查第一个即可。
		return s.VerifySignature(psig, qc.BlockHash)
	}
	return true
	//****

	//var wg sync.WaitGroup
	//var numVerified uint64 = 0
	//for _, psig := range qc.Sigs {
	//	wg.Add(1)
	//	go func(psig PartialSig) { // 实验的时候，模拟即可，开多个gorourine的时间≈一次验证的时间。当节点数很多的时候，goroutine数量太多了，容易打满CPU，所以需要用这个进行模拟。
	//		if s.VerifySignature(psig, qc.BlockHash) {
	//			atomic.AddUint64(&numVerified, 1)
	//		}
	//		wg.Done()
	//	}(psig)
	//}
	//wg.Wait()
	//return numVerified >= uint64(s.conf.QuorumSize)
}

// EvictOld reduces the size of the cache by removing the oldest cached results
func (s *SignatureCache) EvictOld(size int) {
	s.mut.Lock()
	for length := s.cache.Len(); length > size; length-- {
		el := s.cache.Front()
		k := s.cache.Remove(el).(string)
		delete(s.verifiedSignatures, k)
	}
	s.mut.Unlock()
}

// PartialSig is a single replica's signature of a block.
type PartialSig struct {
	ID   config.ReplicaID
	R, S *big.Int
}

func (psig PartialSig) ToBytes() []byte {
	r := psig.R.Bytes()
	s := psig.S.Bytes()
	b := make([]byte, 4, 4+len(r)+len(s))
	binary.LittleEndian.PutUint32(b, uint32(psig.ID))
	b = append(b, r...)
	b = append(b, s...)
	return b
}

// PartialCert is a single replica's certificate for a block.
type PartialCert struct {
	Sig       PartialSig
	BlockHash BlockHash
}

// QuorumCert is a certificate for a block from a quorum of replicas.
type QuorumCert struct {
	Sigs      map[config.ReplicaID]PartialSig
	BlockHash BlockHash
}

func (qc *QuorumCert) ToBytes() []byte {
	b := make([]byte, 0, 32)
	b = append(b, qc.BlockHash[:]...)
	psigs := make([]PartialSig, 0, len(qc.Sigs))
	for _, v := range qc.Sigs {
		i := sort.Search(len(psigs), func(j int) bool {
			return v.ID < psigs[j].ID
		})
		psigs = append(psigs, PartialSig{})
		copy(psigs[i+1:], psigs[i:])
		psigs[i] = v
	}
	for i := range psigs {
		b = append(b, psigs[i].ToBytes()...)
	}
	return b
}

func (qc *QuorumCert) String() string {
	return fmt.Sprintf("QuorumCert{Sigs: %d, Hash: %.8s}", len(qc.Sigs), qc.BlockHash)
}

// AddPartial adds the partial signature to the quorum cert.
func (qc *QuorumCert) AddPartial(cert *PartialCert) error {
	// dont add a cert if there is already a signature from the same replica
	if _, exists := qc.Sigs[cert.Sig.ID]; exists {
		return fmt.Errorf("Attempt to add partial cert from same replica twice")
	}

	if !bytes.Equal(qc.BlockHash[:], cert.BlockHash[:]) { // todo: 实验的时候，这个比较耗时也可以去掉
		return fmt.Errorf("Partial cert hash does not match quorum cert")
	}

	qc.Sigs[cert.Sig.ID] = cert.Sig

	return nil
}

// CreatePartialCert creates a partial cert from a block.
func CreatePartialCert(id config.ReplicaID, privKey *ecdsa.PrivateKey, block *Block) (*PartialCert, error) {
	hash := block.Hash()
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}
	sig := PartialSig{id, r, s}
	return &PartialCert{sig, hash}, nil
}

// VerifyPartialCert will verify a PartialCert from a public key stored in ReplicaConfig。弃用：使用SignatureCache.VerifySignature替代
func VerifyPartialCert(conf *config.ReplicaConfig, cert *PartialCert) bool {
	info, ok := conf.Replicas[cert.Sig.ID]
	if !ok {
		logger.Printf("VerifyPartialSig: got signature from replica whose ID (%d) was not in config.", cert.Sig.ID)
		return false
	}
	return ecdsa.Verify(info.PubKey, cert.BlockHash[:], cert.Sig.R, cert.Sig.S)
}

// CreateQuorumCert creates an empty quorum certificate for a given block
func CreateQuorumCert(block *Block) *QuorumCert {
	return &QuorumCert{BlockHash: block.Hash(), Sigs: make(map[config.ReplicaID]PartialSig)}
}

// VerifyQuorumCert will verify a QuorumCert from public keys stored in ReplicaConfig。弃用：使用SignatureCache.VerifyQuorumCert替代
func VerifyQuorumCert(conf *config.ReplicaConfig, qc *QuorumCert) bool {
	if len(qc.Sigs) < conf.QuorumSize {
		return false
	}
	var wg sync.WaitGroup
	var numVerified uint64 = 0
	for _, psig := range qc.Sigs {
		info, ok := conf.Replicas[psig.ID]
		if !ok {
			logger.Printf("VerifyQuorumSig: got signature from replica whose ID (%d) was not in config.", psig.ID)
		}
		wg.Add(1)
		go func(psig PartialSig) {
			if ecdsa.Verify(info.PubKey, qc.BlockHash[:], psig.R, psig.S) {
				atomic.AddUint64(&numVerified, 1)
			}
			wg.Done()
		}(psig)
	}
	wg.Wait()
	return numVerified >= uint64(conf.QuorumSize)
}
