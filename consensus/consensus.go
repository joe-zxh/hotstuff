package consensus

import (
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
	Blocks   *data.MapStorage
	// BlocksMut sync.Mutex
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

	hs := &HotStuffCore{
		Config:     conf,
		genesis:    genesis,
		bLock:      genesis,
		bExec:      genesis,
		BLeaf:      genesis,
		qcHigh:     qcForGenesis,
		Blocks:     blocks,
		PendingQCs: make(map[data.BlockHash]*data.QuorumCert),
		SigCache:   data.NewSignatureCache(conf),
		cmdCache:   data.NewCommandSet(),
		exec:       make(chan []data.Command, 1),
	}

	hs.waitProposal = sync.NewCond(&hs.Mut)

	hs.Mut.Lock()
	hs.Blocks.Put(genesis)
	hs.Mut.Unlock()

	return hs
}

// UpdateQCHigh updates the qc held by the paceMaker, to the newest qc.
func (hs *HotStuffCore) UpdateQCHigh(qc *data.QuorumCert) bool {
	if !hs.SigCache.VerifyQuorumCert(qc) {
		log.Println("QC not verified!:", qc)
		return false
	}

	logger.Println("UpdateQCHigh")

	newQCHighBlock := hs.GetBlock(qc.BlockHash)

	oldQCHighBlock := hs.GetBlock(hs.qcHigh.BlockHash)

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
	hs.Mut.Lock()
	hs.PutBlock(block)

	qcBlock := hs.GetBlock(block.Justify.BlockHash)

	if block.Height <= hs.vHeight {
		hs.Mut.Unlock()
		logger.Printf("OnReceiveProposal: Block height(%d) less than vHeight(%d)\n", block.Height, hs.vHeight)
		return nil, fmt.Errorf("Block was not accepted")
	}

	safe := false
	if qcBlock.Height > hs.bLock.Height {
		safe = true
	} else {
		logger.Println("OnReceiveProposal: liveness condition failed")
		// check if block extends bLock
		b := block
		for b.Height > hs.bLock.Height+1 {
			b = hs.GetBlock(b.ParentHash)
		}
		if b.ParentHash == hs.bLock.Hash() {
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

	logger.Printf("OnReceiveProposal: Accepted block, height: %d\n", block.Height)
	hs.vHeight = block.Height
	hs.cmdCache.MarkProposed(block.Commands...)

	hs.update(block)

	hs.Mut.Unlock()

	pc, err := hs.SigCache.CreatePartialCert(hs.Config.ID, hs.Config.PrivateKey, block)
	if err != nil {
		panic(err)
		return nil, err
	}
	return pc, nil
}

// should lock
func (hs *HotStuffCore) update(block *data.Block) {
	// block1 = b'', block2 = b', block3 = b
	block1 := hs.GetBlock(block.Justify.BlockHash)
	if block1.Committed {
		return
	}

	logger.Println("PRE COMMIT:", block1)
	// PRE-COMMIT on block1
	hs.UpdateQCHigh(block.Justify)

	block2 := hs.GetBlock(block1.Justify.BlockHash)
	if block2.Committed {
		return
	}

	if block2.Height > hs.bLock.Height {
		hs.bLock = block2 // COMMIT on block2
		logger.Println("COMMIT:", block2)
	}

	block3 := hs.GetBlock(block2.Justify.BlockHash)
	if block3.Committed {
		return
	}

	if block1.ParentHash == block2.Hash() && block2.ParentHash == block3.Hash() {
		logger.Println("DECIDE", block3)
		hs.commit(block3)
		hs.bExec = block3 // DECIDE on block3
	}

	// Free up space by deleting old data
	hs.GCBlocks()
	hs.cmdCache.TrimToLen(hs.Config.BatchSize * 5)
	hs.SigCache.EvictOld(hs.Config.QuorumSize * 5)
}

func (hs *HotStuffCore) commit(block *data.Block) {
	// only called from within update. Thus covered by its mutex lock.
	if hs.bExec.Height < block.Height {
		hs.commit(hs.GetBlock(block.ParentHash)) // todo: 注意genesis
		block.Committed = true
		logger.Println("EXEC", block)
		hs.exec <- block.Commands
	}
}

// CreateProposal creates a new proposal
func (hs *HotStuffCore) CreateProposal() *data.Block {
	batch := hs.cmdCache.GetFirst(hs.Config.BatchSize)
	b := CreateLeaf(hs.BLeaf, batch, hs.qcHigh, hs.BLeaf.Height+1)
	b.Proposer = hs.Config.ID
	hs.PutBlock(b)
	return b
}

// Close frees resources held by HotStuff and closes backend connections
func (hs *HotStuffCore) Close() {
	// hs.cancel()
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

// should lock
func (hs *HotStuffCore) PutBlock(b *data.Block) {
	for {
		if _, ok := hs.Blocks.Get(b.ParentHash); ok { //等到parent到来的時候，才能插入
			hs.Blocks.Put(b)
			hs.waitProposal.Broadcast()
			return
		} else {
			hs.waitProposal.Wait()
		}
	}
}

// should lock
func (hs *HotStuffCore) GetBlock(hash data.BlockHash) *data.Block {
	for {
		if block, ok := hs.Blocks.Get(hash); ok {
			return block
		} else {
			hs.waitProposal.Wait()
		}
	}
}

// should lock
func (hs *HotStuffCore) GCBlocks() {
	hs.Blocks.GarbageCollectBlocks(hs.GetVotedHeight())
}

// should lock
// delete any pending QCs with lower height than bLeaf
func (hs *HotStuffCore) DeletePendings() {
	for k := range hs.PendingQCs {
		if b, ok := hs.Blocks.Get(k); ok {
			if b.Height <= hs.BLeaf.Height {
				delete(hs.PendingQCs, k)
			}
		} /*else { // 不存在的暂时不删吧...
			delete(hs.PendingQCs, k)
		}*/
	}
}
