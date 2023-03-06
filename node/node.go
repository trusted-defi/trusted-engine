package node

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/trusted-defi/trusted-engine/config"
	"github.com/trusted-defi/trusted-engine/core/mempool"
	"github.com/trusted-defi/trusted-engine/log"
	"math/big"
	"path/filepath"
)

const (
	dbfile = "secret.db"
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

func NewNode(nodeconfig config.NodeConfig) *Node {
	n := new(Node)
	n.txpool = mempool.NewTxPool(mempool.DefaultTxPoolConfig, chainConfig, nodeconfig)
	sdbpath := filepath.Join(nodeconfig.NodeDir, dbfile)
	var err error
	if nodeconfig.Generate {
		n.sdb = GenerateDB(sdbpath)
	} else {
		if len(nodeconfig.GivenPrivate) > 0 {
			n.sdb, err = CreateWithHexkey(nodeconfig.GivenPrivate)
			if err != nil {
				log.WithField("err", err).Error("create secret db with key failed")
			} else {
				SaveDb(n.sdb, sdbpath)
			}
		} else {
			n.sdb = LoadDb(sdbpath)
		}
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
