package chainclient

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/trusted-defi/trusted-engine/config"
	"math/big"
	"testing"
)

func TestChainClient_CurrentBlock(t *testing.T) {
	client, err := NewChainClient(config.GetConfig())
	if err != nil {
		t.Error("new chain client failed", err)
	}
	if b, err := client.CurrentBlock(); err != nil {
		fmt.Printf("get currentBlock is empty\n")
	} else {
		fmt.Printf("get currentBlock, number is %d\n", b.Number().Uint64())
	}
}

func TestChainClient_GetBalance(t *testing.T) {
	client, err := NewChainClient(config.GetConfig())
	if err != nil {
		t.Error("new chain client failed", err)
	}
	addr := common.HexToAddress("0x5c27dd97ddf34588006740a2c2665f5dba3b76c8")
	b := client.GetBalance(addr)
	if b == nil {
		fmt.Printf("get balance is empty\n")
	} else {
		fmt.Printf("get balance is %s\n", b.Text(10))
	}
}

func TestChainClient_NonceAt(t *testing.T) {
	client, err := NewChainClient(config.GetConfig())
	if err != nil {
		t.Error("new chain client failed", err)
	}
	addr := common.HexToAddress("0x5c27dd97ddf34588006740a2c2665f5dba3b76c8")
	b := client.NonceAt(addr)
	fmt.Printf("get nonce is %d\n", b)
}

func TestChainClient_NonceAtHeight(t *testing.T) {
	client, err := NewChainClient(config.GetConfig())
	if err != nil {
		t.Error("new chain client failed", err)
	}
	addr := common.HexToAddress("0x5c27dd97ddf34588006740a2c2665f5dba3b76c8")
	b := client.NonceAtHeight(addr, big.NewInt(100))
	fmt.Printf("get nonce is %d\n", b)
}

func TestChainClient_SubscribeChainHeadEvent(t *testing.T) {
	client, err := NewChainClient(config.GetConfig())
	if err != nil {
		t.Error("new chain client failed", err)
	}
	client.Start()
	ch := make(chan core.ChainHeadEvent, 10)
	sub := client.SubscribeChainHeadEvent(ch)
	for i := 0; i < 5; i++ {
		select {
		case <-sub.Err():
			fmt.Printf("subscribe error %s\n", err.Error())
		case event, ok := <-ch:
			if !ok {
				fmt.Printf("get event failed\n")
				return
			}
			fmt.Printf("got new head.number = %d\n", event.Block.Number().Uint64())

		}
	}
}
