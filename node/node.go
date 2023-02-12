package node

import "github.com/trusted-defi/trusted-engine/core/mempool"

type Node struct {
	txpool *mempool.TxPool
}

func NewNode() *Node {
	n := new(Node)
	n.txpool = mempool.NewTxPool(mempool.DefaultTxPoolConfig, nil)

	return n
}

func (n *Node) TxPool() *mempool.TxPool {
	return n.txpool
}

func (n *Node) IsReady() bool {
	return n.txpool.IsReady()
}
