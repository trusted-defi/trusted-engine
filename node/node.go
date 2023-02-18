package node

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/trusted-defi/trusted-engine/core/mempool"
	"math/big"
	"path/filepath"
)

const (
	NodeDir = "nodedata"
	dbfile  = "secret.db"
)

type Node struct {
	txpool *mempool.TxPool
	sdb    *SecretDb
}

func init() {
	// todo: set chainid with chain client.
	updateConfig()
}

var (
	chainConfig = params.MainnetChainConfig
)

func updateConfig() {
	chainConfig.ChainID = big.NewInt(1024)
}

func NewNode(generate bool) *Node {
	n := new(Node)
	n.txpool = mempool.NewTxPool(mempool.DefaultTxPoolConfig, chainConfig, NodeDir)
	sdbpath := filepath.Join(NodeDir, dbfile)
	if generate {
		n.sdb = GenerateDB(sdbpath)
	} else {
		n.sdb = LoadDb(sdbpath)
	}

	return n
}

func (n *Node) TxPool() *mempool.TxPool {
	return n.txpool
}

func (n *Node) IsReady() bool {
	return n.txpool.IsReady()
}

func (n *Node) SetPrivk(hexk string) error {
	if sdb, err := CreateWithHexkey(hexk); err != nil {
		return err
	} else {
		n.sdb = sdb
	}
	return nil
}

func (n *Node) GetSecretDB() *SecretDb {
	return n.sdb
}
