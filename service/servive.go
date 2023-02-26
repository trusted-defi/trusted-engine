package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	corecmn "github.com/trusted-defi/trusted-engine/core/common"
	"github.com/trusted-defi/trusted-engine/core/cryptor"
	"github.com/trusted-defi/trusted-engine/core/mempool"
	"github.com/trusted-defi/trusted-engine/node"
	trusted "github.com/trusted-defi/trusted-engine/protocol/generate/trusted/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/big"
	"net"
)

var log = logrus.WithField("prefix", "service")

type TrustedService struct {
	n *node.Node
	trusted.UnimplementedTrustedServiceServer
}

func parseTxStatus(status []mempool.TxStatus) []uint32 {
	s := make([]uint32, 0, len(status))
	for _, st := range status {
		s = append(s, uint32(st))
	}
	return s
}
func parseHashs(hashs_data [][]byte) []common.Hash {
	hashs := make([]common.Hash, 0, len(hashs_data))
	for _, hash_data := range hashs_data {
		h := common.BytesToHash(hash_data)
		hashs = append(hashs, h)
	}
	return hashs
}

func parseErrs(errs []error) []string {
	strs := make([]string, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			strs = append(strs, "")
		} else {
			strs = append(strs, err.Error())
		}
	}
	return strs
}

func toBigInt(data []byte) *big.Int {
	return new(big.Int).SetBytes(data)
}

func parseListToTransactions(list *trusted.TransactionList) []*types.Transaction {
	txs := make([]*types.Transaction, 0, len(list.Txs))
	for _, txdata := range list.Txs {
		tx := corecmn.ParseTxData(txdata)
		txs = append(txs, tx)
	}
	return txs
}

func sliceToList(account common.Address, txs types.Transactions) *trusted.AccountTransactionList {
	txlist := &trusted.AccountTransactionList{
		Address: account.Bytes(),
		TxList: &trusted.TransactionList{
			Txs: make([][]byte, 0, len(txs)),
		},
	}
	for _, tx := range txs {
		if encode, err := tx.MarshalBinary(); err != nil {
			//
		} else {
			txlist.TxList.Txs = append(txlist.TxList.Txs, encode)
		}
	}
	return txlist
}

func maptxToList(account_txs map[common.Address]types.Transactions) []*trusted.AccountTransactionList {
	txlists := make([]*trusted.AccountTransactionList, 0, len(account_txs))
	for addr, txs := range account_txs {
		txlist := sliceToList(addr, txs)
		txlists = append(txlists, txlist)
	}
	return txlists
}

func parseAddrsToBytes(accounts []common.Address) [][]byte {
	account_data := make([][]byte, 0, len(accounts))
	for _, account := range accounts {
		data := account.Bytes()
		account_data = append(account_data, data)
	}
	return account_data
}

func (s *TrustedService) parseCryptedTxsTransactions(request *trusted.AddTrustedTxsRequest) ([]*types.Transaction, []error) {
	txs := make([]*types.Transaction, len(request.CtyptedTxs))
	errs := make([]error, len(request.CtyptedTxs))
	for i, cryptedTx := range request.CtyptedTxs {

		txdata, err := s.decrypt(cryptedTx)
		if err != nil {
			errs[i] = err
			continue
		}
		tx := new(types.Transaction)
		err = tx.UnmarshalBinary(txdata)
		if err != nil {
			errs[i] = err
			continue
		}
		txs[i] = tx

	}
	return txs, errs
}

func (s *TrustedService) ServiceReady(ctx context.Context, req *emptypb.Empty) (*trusted.ServiceReadyResponse, error) {
	res := new(trusted.ServiceReadyResponse)
	res.Ready = s.n.IsReady()
	return res, nil
}

func (s *TrustedService) PoolSetPrice(ctx context.Context, req *trusted.SetPriceRequest) (*emptypb.Empty, error) {
	s.n.TxPool().SetGasPrice(toBigInt(req.Price))
	return nil, nil
}

func (s *TrustedService) PoolGasPrice(ctx context.Context, req *emptypb.Empty) (*trusted.GasPriceResponse, error) {
	res := new(trusted.GasPriceResponse)
	res.Price = s.n.TxPool().GasPrice().Bytes()
	return res, nil
}

func (s *TrustedService) PendingNonce(ctx context.Context, req *trusted.PendingNonceRequest) (*trusted.PendingNonceResponse, error) {
	addr := common.BytesToAddress(req.Address)
	nonce := s.n.TxPool().Nonce(addr)
	res := new(trusted.PendingNonceResponse)
	res.Nonce = nonce

	return res, nil
}

func (s *TrustedService) PoolStat(ctx context.Context, req *emptypb.Empty) (*trusted.PoolStatResponse, error) {
	res := new(trusted.PoolStatResponse)
	pending, queue := s.n.TxPool().Stats()
	//log.WithField("pending", pending).WithField("queue", queue).Info("txpool stat")
	log.WithField("pending", pending).Info("txpool stat")
	res.Pending = uint64(pending)
	res.Queue = uint64(queue)
	return res, nil
}

func (s *TrustedService) PoolContent(ctx context.Context, req *trusted.PoolContentRequest) (*trusted.PoolContentResponse, error) {
	pendings, queue := s.n.TxPool().Content()
	res := new(trusted.PoolContentResponse)
	res.PendingList = maptxToList(pendings)
	res.QueueList = maptxToList(queue)

	return res, nil
}
func (s *TrustedService) PoolContentFrom(ctx context.Context, req *trusted.PoolContentRequest) (*trusted.PoolContentResponse, error) {
	from := common.BytesToAddress(req.Address)
	res := new(trusted.PoolContentResponse)
	pendings, queue := s.n.TxPool().ContentFrom(from)
	res.PendingList = []*trusted.AccountTransactionList{sliceToList(from, pendings)}
	res.QueueList = []*trusted.AccountTransactionList{sliceToList(from, queue)}
	return res, nil
}

func (s *TrustedService) PoolPending(ctx context.Context, req *emptypb.Empty) (*trusted.PoolPendingResponse, error) {
	pending := s.n.TxPool().Pending(false)
	list := maptxToList(pending)
	res := new(trusted.PoolPendingResponse)
	res.PendingList = list
	return res, nil
}

func (s *TrustedService) PoolLocals(ctx context.Context, req *emptypb.Empty) (*trusted.PoolLocalsResponse, error) {
	locals := s.n.TxPool().Locals()
	res := new(trusted.PoolLocalsResponse)
	res.AddressList = parseAddrsToBytes(locals)
	return res, nil
}

func (s *TrustedService) AddLocalsTx(ctx context.Context, req *trusted.AddTxsRequest) (*trusted.AddTxsResponse, error) {
	txs := parseListToTransactions(req.TxList)
	errs := s.n.TxPool().AddLocals(txs)
	res := new(trusted.AddTxsResponse)
	res.Errors = parseErrs(errs)
	return res, nil
}

func (s *TrustedService) AddRemoteTx(ctx context.Context, req *trusted.AddTxsRequest) (*trusted.AddTxsResponse, error) {
	txs := parseListToTransactions(req.TxList)
	errs := s.n.TxPool().AddRemotes(txs)
	res := new(trusted.AddTxsResponse)
	res.Errors = parseErrs(errs)
	return res, nil
}

func (s *TrustedService) TxStatus(ctx context.Context, req *trusted.TxStatusRequest) (*trusted.TxStatusResponse, error) {
	txhashs := parseHashs(req.TxHashs)
	txstatus := s.n.TxPool().Status(txhashs)
	res := new(trusted.TxStatusResponse)
	res.TxStatus = parseTxStatus(txstatus)
	return res, nil
}

func (s *TrustedService) TxGet(ctx context.Context, req *trusted.TxGetRequest) (*trusted.TxGetResponse, error) {
	txhash := common.BytesToHash(req.TxHash)
	tx := s.n.TxPool().Get(txhash)
	res := new(trusted.TxGetResponse)
	res.Tx, _ = tx.MarshalBinary()
	return res, nil
}
func (s *TrustedService) TxHas(ctx context.Context, req *trusted.TxHasRequest) (*trusted.TxHasResponse, error) {
	txhash := common.BytesToHash(req.TxHash)
	res := new(trusted.TxHasResponse)
	res.Has = s.n.TxPool().Has(txhash)
	return res, nil
}

func (s *TrustedService) SubscribeNewTransaction(req *trusted.SubscribeNewTxRequest, server trusted.TrustedService_SubscribeNewTransactionServer) error {
	eventCh := make(chan core.NewTxsEvent)
	sub := s.n.TxPool().SubscribeNewTxsEvent(eventCh)
	var bcontinue = true
	for bcontinue {
		select {
		case err := <-sub.Err():
			log.Println("subscribe failed", err)
			return err
		case event := <-eventCh:
			var crypted_tx_data = make([][]byte, 0)
			for _, tx := range event.Txs {
				txdata, err := tx.MarshalBinary()
				if err != nil {
					continue
				}
				crypted, err := s.crypt(txdata)
				if err != nil {
					log.Println("crypt tx data failed", err)
					continue
				}
				crypted_tx_data = append(crypted_tx_data, crypted)
			}
			response := trusted.SubscribeNewTxResponse{
				CryptedNewTx: crypted_tx_data,
			}
			server.Send(&response)
		}
	}
	return nil
}

func (s *TrustedService) Crypt(ctx context.Context, req *trusted.CryptRequest) (*trusted.CryptResponse, error) {
	var err error
	res := new(trusted.CryptResponse)
	res.Crypted, err = s.crypt(req.GetData())
	if err != nil {
		return nil, err
	}
	return res, nil
}

func generateTxAsset(tx *types.Transaction) []byte {
	// todo: change generate method to sgx.
	return []byte{1}
}

func (s *TrustedService) crypt(data []byte) ([]byte, error) {
	sdb := s.n.GetSecretDB()
	if sdb != nil {
		return cryptor.Encrypt(data, sdb.PublicKey())
	} else {
		return nil, errors.New("public key not found")
	}
}

func (s *TrustedService) decrypt(data []byte) ([]byte, error) {
	sdb := s.n.GetSecretDB()
	if sdb != nil {
		return cryptor.Decrypt(data, sdb.PrivateKey())
	} else {
		return nil, errors.New("private key not found")
	}
}

func (s *TrustedService) AddLocalTrustedTxs(ctx context.Context, req *trusted.AddTrustedTxsRequest) (*trusted.AddTrustedTxsResponse, error) {
	res := new(trusted.AddTrustedTxsResponse)
	res.Results = make([]*trusted.AddTrustedTxResult, len(req.CtyptedTxs))

	txs, parsedErrs := s.parseCryptedTxsTransactions(req)
	addTxErrs := s.n.TxPool().AddLocals(txs)
	for i, tx := range txs {
		result := new(trusted.AddTrustedTxResult)
		parseErr := parsedErrs[i]
		addTxErr := addTxErrs[i]
		if parseErr != nil {
			result.Hash = common.Hash{}.Bytes()
			result.Asset = make([]byte, 0)
			result.Error = parseErr.Error()
		} else if addTxErr != nil {
			result.Hash = common.Hash{}.Bytes()
			result.Asset = make([]byte, 0)
			result.Error = addTxErr.Error()
		} else {
			result.Hash = tx.Hash().Bytes()
			result.Asset = generateTxAsset(tx)
			result.Error = ""
		}
		res.Results[i] = result
	}

	return res, nil
}

func (s *TrustedService) AddRemoteTrustedTx(ctx context.Context, req *trusted.AddTrustedTxsRequest) (*trusted.AddTrustedTxsResponse, error) {
	res := new(trusted.AddTrustedTxsResponse)
	res.Results = make([]*trusted.AddTrustedTxResult, len(req.CtyptedTxs))

	txs, parsedErrs := s.parseCryptedTxsTransactions(req)
	addTxErrs := s.n.TxPool().AddRemotes(txs)
	for i, tx := range txs {
		result := new(trusted.AddTrustedTxResult)
		parseErr := parsedErrs[i]
		addTxErr := addTxErrs[i]
		if parseErr != nil {
			result.Hash = common.Hash{}.Bytes()
			result.Asset = make([]byte, 0)
			result.Error = parseErr.Error()
		} else if addTxErr != nil {
			result.Hash = common.Hash{}.Bytes()
			result.Asset = make([]byte, 0)
			result.Error = addTxErr.Error()
		} else {
			result.Hash = tx.Hash().Bytes()
			result.Asset = generateTxAsset(tx)
			result.Error = ""
		}
		res.Results[i] = result
	}

	return res, nil
}

func (s *TrustedService) FillBlock(ctx context.Context, req *trusted.FillBlockRequest) (*trusted.FillBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FillBlock not implemented")
}

func RegisterService(server *grpc.Server, n *node.Node) {
	s := new(TrustedService)
	s.n = n
	trusted.RegisterTrustedServiceServer(server, s)
}

func StartTrustedService(n *node.Node) {
	lis, err := net.Listen("tcp", ":3802")
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
