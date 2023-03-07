package smanager

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/edgelesssys/ego/enclave"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/trusted-defi/trusted-engine/log"
	"strings"
	"sync"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
)

type progress struct {
	randomA, randomB, randomC    []byte
	randomAR, randomBR, randomCR []byte
}

type WatchKeyHandler func([]byte)

type KeyManager struct {
	privk    *ecies.PrivateKey
	sessions map[string]*progress
	mux      sync.Mutex
	watchers []WatchKeyHandler
}

func NewKeyManager(pk *ecies.PrivateKey) *KeyManager {
	return &KeyManager{
		privk:    pk,
		sessions: make(map[string]*progress),
		watchers: make([]WatchKeyHandler, 0),
	}
}

func (t *KeyManager) AddKeyWatcher(handler WatchKeyHandler) {
	t.mux.Lock()
	defer t.mux.Unlock()
	if handler != nil {
		t.watchers = append(t.watchers, handler)
	}
}

// CheckSecretKey check secretkey already exist or not.
func (t *KeyManager) CheckSecretKey() bool {
	t.mux.Lock()
	defer t.mux.Unlock()

	return t.privk != nil
}

// GetAuthData generate a remote report at begin of a auth-verify process.
func (t *KeyManager) GetAuthData(peerId string) ([]byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	var pro *progress
	if p, exist := t.sessions[strings.ToLower(peerId)]; exist {
		pro = p
	} else {
		pro = new(progress)
		pro.randomA = GenRandom()
		t.sessions[strings.ToLower(peerId)] = pro
	}
	report, err := enclave.GetRemoteReport(pro.randomA)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// VerifyAuth verify auth data received from remote peer
func (t *KeyManager) VerifyAuth(authData []byte, peerId string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return ErrInvalidOperation
	} else {
		report, err := enclave.VerifyRemoteReport(authData)
		if err != nil {
			return err
		}
		p.randomAR = common.CopyBytes(report.Data)
	}
	return nil
}

// GetVerifyData generate a remote report used to verify remote peer..
func (t *KeyManager) GetVerifyData(peerId string) ([]byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return nil, ErrInvalidOperation
	} else {
		p.randomB = GenRandom()
		data := append(p.randomAR, p.randomB...)
		report, err := enclave.GetRemoteReport(data)
		if err != nil {
			return nil, err
		}
		return report, nil
	}
}

// VerifyRemoteVerify verify remote verify-data received from remote peer..
func (t *KeyManager) VerifyRemoteVerify(verifyData []byte, peerId string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return ErrInvalidOperation
	} else {
		report, err := enclave.VerifyRemoteReport(verifyData)
		if err != nil {
			return err
		}
		if len(report.Data) != 64 {
			return ErrInvalidOperation
		}
		if bytes.Compare(p.randomA, report.Data[:32]) != 0 {
			return ErrInvalidOperation
		}
		p.randomBR = common.CopyBytes(report.Data[32:])
	}
	return nil
}

// GetRequestKeyData generate a remote report used to request secret key.
func (t *KeyManager) GetRequestKeyData(peerId string) ([]byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return nil, ErrInvalidOperation
	} else {
		p.randomC = GenRandom()
		data := append(p.randomBR, p.randomC...)
		report, err := enclave.GetRemoteReport(data)
		if err != nil {
			return nil, err
		}
		return report, nil
	}
}

// VerifyRequestKeyData verify remote verify-data received from remote peer..
func (t *KeyManager) VerifyRequestKeyData(request []byte, peerId string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return ErrInvalidOperation
	} else {
		report, err := enclave.VerifyRemoteReport(request)
		if err != nil {
			return err
		}
		if len(report.Data) != 64 {
			return ErrInvalidOperation
		}
		if bytes.Compare(p.randomB, report.Data[:32]) != 0 {
			return ErrInvalidOperation
		}
		p.randomCR = common.CopyBytes(report.Data[32:])
	}
	return nil
}

// GetResponseKeyData generate a remote report used to request secret key.
func (t *KeyManager) GetResponseKeyData(peerId string) ([]byte, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return nil, ErrInvalidOperation
	} else {
		if t.privk == nil {
			return nil, ErrInvalidOperation
		}
		hexk := ethcmn.Bytes2Hex(crypto.FromECDSA(t.privk.ExportECDSA()))
		key := ethcmn.FromHex(hexk)
		data := append(p.randomCR, key...)
		report, err := enclave.GetRemoteReport(data)
		if err != nil {
			return nil, err
		}
		return report, nil
	}
}

// VerifyResponseKey verify remote verify-data received from remote peer..
func (t *KeyManager) VerifyResponseKey(response []byte, peerId string) error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return ErrInvalidOperation
	} else {
		report, err := enclave.VerifyRemoteReport(response)
		if err != nil {
			return err
		}
		if len(report.Data) != 64 {
			return ErrInvalidOperation
		}
		if bytes.Compare(p.randomC, report.Data[:32]) != 0 {
			return ErrInvalidOperation
		}
		key := report.Data[32:]
		pk, err := crypto.HexToECDSA(ethcmn.Bytes2Hex(key))
		if err != nil {
			return err
		}
		if t.privk != nil {
			log.Warn("key manager have been store private key, skip new key")
			return nil
		}
		t.privk = ecies.ImportECDSA(pk)
		for _, handler := range t.watchers {
			handler(common.CopyBytes(key))
		}
		log.Info("key manager got private key")
	}
	return nil
}

func GenRandom() []byte {
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return crypto.FromECDSA(pk)
}
