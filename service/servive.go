package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	corecmn "github.com/trusted-defi/trusted-engine/core/common"
	"github.com/trusted-defi/trusted-engine/node"
	trusted "github.com/trusted-defi/trusted-engine/protocol/generate/trusted/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type TrustedService struct {
	n *node.Node
	trusted.UnimplementedTrustedServiceServer
}

func (s *TrustedService) HealthCheck(ctx context.Context, req *emptypb.Empty) (*trusted.HealthCheckResponse, error) {
	res := new(trusted.HealthCheckResponse)
	res.Status = true
	return res, nil
}
func (s *TrustedService) AddTx(ctx context.Context, req *trusted.AddTxRequest) (*trusted.AddTxResponse, error) {
	res := new(trusted.AddTxResponse)
	tx := corecmn.ParseTxData(req.Txdata)
	if tx != nil {
		err := s.n.TxPool().AddLocal(tx)
		if err != nil {
			res.Error = err.Error()
		}
		res.Txhash = tx.Hash().Bytes()
	} else {
		res.Error = errors.New("invalid tx data").Error()
	}

	return res, nil
}

func (s *TrustedService) Status(ctx context.Context, req *trusted.StatusRequest) (*trusted.StatusResponse, error) {
	pending, queue := s.n.TxPool().Stats()
	res := new(trusted.StatusResponse)
	res.Pending = uint64(pending)
	res.Queue = uint64(queue)
	return res, nil
}
func (s *TrustedService) Reset(ctx context.Context, req *trusted.ResetRequest) (*trusted.ResetResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}
func (s *TrustedService) Pending(ctx context.Context, req *trusted.PendingRequest) (*trusted.PendingResponse, error) {
	alltxs := s.n.TxPool().Pending(false)
	res := new(trusted.PendingResponse)
	res.Txs = make([][]byte, 0)
	for _, txs := range alltxs {
		for _, tx := range txs {
			d, _ := tx.MarshalBinary()
			res.Txs = append(res.Txs, d)
		}
	}
	return res, nil
}

func (s *TrustedService) Crypt(ctx context.Context, req *trusted.CryptRequest) (*trusted.CryptResponse, error) {
	var err error
	res := new(trusted.CryptResponse)
	res.Crypted, err = crypt(req.GetData())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func generateTxAsset(tx *types.Transaction) []byte {
	// todo: change generate method to sgx.
	return []byte{1}
}

func crypt(data []byte) ([]byte, error) {
	// todo: change crypt method to sgx.
	return data, nil
}

func decrypt(data []byte) ([]byte, error) {
	// todo: change decrypt method to sgx.
	return data, nil
}
func (s *TrustedService) AddLocalTrustedTx(ctx context.Context, req *trusted.AddTrustedTxRequest) (*trusted.AddTrustedTxResponse, error) {
	res := new(trusted.AddTrustedTxResponse)
	tx := new(types.Transaction)
	txdata, err := decrypt(req.GetCtyptedTx())
	if err != nil {
		return nil, err
	}
	err = tx.UnmarshalBinary(txdata)
	if err != nil {
		res.Hash = common.Hash{}.Bytes()
		res.Asset = make([]byte, 0)
		return res, err
	}
	// todo: add check tx and add local tx to txpool.

	res.Hash = tx.Hash().Bytes()
	res.Asset = generateTxAsset(tx)
	return res, nil
}

func (s *TrustedService) AddRemoteTrustedTx(ctx context.Context, req *trusted.AddTrustedTxRequest) (*trusted.AddTrustedTxResponse, error) {
	res := new(trusted.AddTrustedTxResponse)
	tx := new(types.Transaction)
	txdata, err := decrypt(req.GetCtyptedTx())
	if err != nil {
		return nil, err
	}
	err = tx.UnmarshalBinary(txdata)
	if err != nil {
		res.Hash = common.Hash{}.Bytes()
		res.Asset = make([]byte, 0)
		return res, err
	}
	// todo: add check tx and add remote tx to txpool.

	res.Hash = tx.Hash().Bytes()
	res.Asset = generateTxAsset(tx)
	return res, nil
}

func RegisterService(server *grpc.Server, n *node.Node) {
	s := new(TrustedService)
	s.n = n
	trusted.RegisterTrustedServiceServer(server, s)
}

func StartTrustedService(n *node.Node) {
	lis, err := net.Listen("tcp", ":38000")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	RegisterService(s, n)

	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
