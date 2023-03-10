package chainclient

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/trusted-defi/trusted-engine/config"
	corecmn "github.com/trusted-defi/trusted-engine/core/common"
	"github.com/trusted-defi/trusted-engine/log"
	trusted "github.com/trusted-defi/trusted-engine/protocol/generate/trusted/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
	"time"
)

type ChainClient struct {
	cclient trusted.ChainServiceClient

	chainHeadFeed event.Feed
	scope         event.SubscriptionScope
	quit          chan struct{}
	ctx           context.Context
}

func NewChainClient(nodeconfig config.NodeConfig) (*ChainClient, error) {
	c, err := grpc.Dial(nodeconfig.ChainServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("dial server failed")
	}

	log.Info("grpc connected")
	client := new(ChainClient)
	client.ctx = context.Background()
	client.cclient = trusted.NewChainServiceClient(c)
	client.quit = make(chan struct{})
	client.Start()

	return client, nil
}

func (client *ChainClient) Start() {
	go client.loop()
}

func (client *ChainClient) CurrentBlock() (*types.Block, error) {
	latest, err := client.cclient.CurrentBlock(client.ctx, new(trusted.CurrentBlockRequest), grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get current block failed", "err", err)
		return nil, err
	}
	return corecmn.ParseBlockData(latest.BlockData), nil
}

func (client *ChainClient) GetBlock(hash common.Hash, number uint64) *types.Block {
	req := new(trusted.BlockRequest)
	req.BlockHash = hash.Bytes()
	req.BlockNum = number
	block, err := client.cclient.GetBlock(client.ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get current block failed", "err", err)
		return nil
	}
	return corecmn.ParseBlockData(block.BlockData)
}

func (client *ChainClient) GetBalance(addr common.Address) *big.Int {
	req := new(trusted.BalanceRequest)
	req.Address = addr.Bytes()
	balance, err := client.cclient.GetBalance(client.ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get balance failed", "err", err)
		return big.NewInt(0)
	}
	return corecmn.ParseBalance(balance)
}

func (client *ChainClient) NonceAtHeight(addr common.Address, height *big.Int) uint64 {
	req := new(trusted.NonceRequest)
	req.Address = addr.Bytes()
	req.BlockNum = height.Bytes()
	nonce, err := client.cclient.GetNonce(client.ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get balance failed", "err", err)
		return 0
	}
	return corecmn.ParseNonce(nonce)
}

func (client *ChainClient) NonceAt(addr common.Address) uint64 {
	req := new(trusted.NonceRequest)
	req.Address = addr.Bytes()
	nonce, err := client.cclient.GetNonce(client.ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get balance failed", "err", err)
		return 0
	}
	return corecmn.ParseNonce(nonce)
}

func (client *ChainClient) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return client.scope.Track(client.chainHeadFeed.Subscribe(ch))
}

func (client *ChainClient) loop() {
	for {
		subsucceed := false
		sub, err := client.cclient.ChainHeadEvent(client.ctx, new(trusted.ChainHeadEventRequest))
		if err != nil {
			log.Error("chain head event subscribe failed", "err", err)
			time.Sleep(time.Second)
			continue
		}
		subsucceed = true
		for subsucceed {
			select {
			case <-client.quit:
				log.Info("chain client quit")
				return
			default:

			}

			if res, err := sub.Recv(); err != nil {
				log.Info("chain head event receive failed", "err", err)
				subsucceed = false
			} else {
				client.chainHeadFeed.Send(core.ChainHeadEvent{
					Block: corecmn.ParseBlockData(res.BlockData),
				})
			}
		}
	}
}
