package node

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/trusted-defi/trusted-engine/core/mempool"
	"math/big"
)

type Node struct {
	txpool *mempool.TxPool
}

func init() {
	updateConfig()
}

var (
	chainConfig = params.MainnetChainConfig
)

func updateConfig() {
	chainConfig.ChainID = big.NewInt(1024)
}

func NewNode() *Node {
	n := new(Node)
	n.txpool = mempool.NewTxPool(mempool.DefaultTxPoolConfig, chainConfig)

	return n
}

func (n *Node) TxPool() *mempool.TxPool {
	return n.txpool
}

func (n *Node) IsReady() bool {
	return n.txpool.IsReady()
}
