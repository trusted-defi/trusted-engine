package blockfill

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"sync"
)

type BlockFillRecord struct {
	BlockHash    common.Hash `json:"block-hash"`
	MatchTxCount int         `json:"match-tx-count"`
	ProofRoot    common.Hash `json:"proof-root"`
	BlockRoot    common.Hash `json:"block-root"`
}

type BlockProofParam struct {
	ParentHash common.Hash
	BlockTime  uint64
}

type BlockProof struct {
	givenTxs []*types.Transaction
	txsroot  common.Hash
}

type BlockFiller struct {
	mux   sync.Mutex
	proof map[common.Hash]*BlockProof
}

func NewBlockFiller(datapath string) *BlockFiller {
	return &BlockFiller{
		proof: make(map[common.Hash]*BlockProof),
	}
}

func (b *BlockFiller) SetBlockProof(param BlockProofParam, txs []*types.Transaction) {
	blockparam := fmt.Sprintf("%s%v", param.ParentHash.String(), param.BlockTime)
	bid := crypto.Keccak256Hash([]byte(blockparam))
	proof := &BlockProof{
		givenTxs: txs,
	}
	proof.txsroot = types.DeriveSha(types.Transactions(txs), trie.NewStackTrie(nil))
	b.proof[bid] = proof
}

func (b *BlockFiller) VerifyBlock(block *types.Block) {
	blockparam := fmt.Sprintf("%s%v", block.ParentHash().String(), block.Time())
	bid := crypto.Keccak256Hash([]byte(blockparam))
	proof, exist := b.proof[bid]
	if !exist {
		return
	}
	txmap := make(map[common.Hash]struct{})
	for _, tx := range proof.givenTxs {
		txmap[tx.Hash()] = struct{}{}
	}
	matchCount := 0
	for _, tx := range block.Transactions() {
		if _, exist := txmap[tx.Hash()]; exist {
			matchCount++
		}
	}

	record := &BlockFillRecord{
		BlockHash:    block.Hash(),
		BlockRoot:    block.TxHash(),
		ProofRoot:    proof.txsroot,
		MatchTxCount: matchCount,
	}
	b.saveRecord(record)
}

func (b *BlockFiller) saveRecord(record *BlockFillRecord) {
	// todo: save record to db.
}
