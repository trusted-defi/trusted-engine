package mempool

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/trusted-defi/trusted-engine/config"
	corecmn "github.com/trusted-defi/trusted-engine/core/common"
	trusted "github.com/trusted-defi/trusted-engine/protocol/generate/trusted/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
)

type ChainClient struct {
	cclient trusted.ChainServiceClient

	chainHeadFeed event.Feed
	scope         event.SubscriptionScope

	ctx context.Context
}

func NewChainClient(conf *config.Config) (*ChainClient, error){
	c,err := grpc.Dial(conf.ChainServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("dial server failed")
	}
	client := new(ChainClient)
	client.ctx = context.Background()
	client.cclient = trusted.NewChainServiceClient(c)
	go client.loop()
	return client, nil
}

func (client *ChainClient) CurrentBlock() *types.Block {
	latest, err := client.cclient.CurrentBlock(client.ctx, new(trusted.CurrentBlockRequest), grpc.EmptyCallOption{})
	if err != nil {
		log.Error("get current block failed", "err", err)
		return nil
	}
	return corecmn.ParseBlockData(latest.BlockData)
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

func (client *ChainClient) loop () {
	for {
		sub, _ := client.cclient.ChainHeadEvent(client.ctx, new(trusted.ChainHeadEventRequest))
		for {
			if res, err := sub.Recv(); err != nil {
				log.Info("chain head event receive failed", "err", err)
				break
			} else {
				client.chainHeadFeed.Send(core.ChainHeadEvent{
					Block: corecmn.ParseBlockData(res.BlockData),
				})
			}
		}
	}
}