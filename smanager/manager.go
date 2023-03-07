package smanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"github.com/edgelesssys/ego/enclave"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"strings"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
)

type progress struct {
	randomA, randomB, randomC    []byte
	randomAR, randomBR, randomCR []byte
}

type KeyManager struct {
	privk    *ecies.PrivateKey
	sessions map[string]*progress
}

// CheckSecretKey check secretkey already exist or not.
func (t *KeyManager) CheckSecretKey() bool {
	return t.privk != nil
}

// GetAuthData generate a remote report at begin of a auth-verify process.
func (t *KeyManager) GetAuthData(peerId string) ([]byte, error) {
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
	if p, exist := t.sessions[strings.ToLower(peerId)]; !exist {
		return ErrInvalidOperation
	} else {
		report, err := enclave.VerifyRemoteReport(verifyData)
		if err != nil {
			return err
		}

		p.randomAR = common.CopyBytes(report.Data)
	}
	return nil
}

// GetRequestKeyData generate a remote report used to request secret key.
func (t *KeyManager) GetRequestKeyData(peerId string) ([]byte, error) {
	//todo: implement
	return nil, ErrClientNotReady
}

// VerifyRequestKeyData verify remote verify-data received from remote peer..
func (t *KeyManager) VerifyRequestKeyData(request []byte, peerId string) error {
	//todo: implement
	return ErrClientNotReady
}

// GetResponseKeyData generate a remote report used to request secret key.
func (t *KeyManager) GetResponseKeyData(peerId string) ([]byte, error) {
	//todo: implement
	return nil, ErrClientNotReady
}

// VerifyResponseKey verify remote verify-data received from remote peer..
func (t *KeyManager) VerifyResponseKey(response []byte, peerId string) error {
	//todo: implement
	return ErrClientNotReady
}

func GenRandom() []byte {
	pk, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return crypto.FromECDSA(pk)
}
