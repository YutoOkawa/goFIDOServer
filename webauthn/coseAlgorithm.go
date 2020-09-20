package webauthn

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor"
)

type PublicKey struct {
	KeyType int64 `cbor:"1,keyasint" json:"kty"`
	Alg     int64 `cbor:"3,keyasint" json:"alg"`
}

type EC2PublicKey struct {
	PublicKey
	Crv int64  `cbor:"-1,keyasint,omitempty" json:"crv"`
	X   []byte `cbor:"-2,keyasint,omitempty" json:"x"`
	Y   []byte `cbor:"-3,keyasint,omitempty" json:"y"`
}

type ECDSASignature struct {
	R *big.Int
	S *big.Int
}

func parsePublicKey(pubkeyBytes []byte) (interface{}, error) {
	var pubKey PublicKey
	err := cbor.Unmarshal(pubkeyBytes, &pubKey)
	if err != nil {
		return nil, err
	}

	switch pubKey.KeyType {
	case 2:
		var ePubKey EC2PublicKey
		err := cbor.Unmarshal(pubkeyBytes, &ePubKey)
		if err != nil {
			return nil, err
		}
		ePubKey.PublicKey = pubKey
		return ePubKey, nil
	default:
		return nil, fmt.Errorf("Not Implemented Error")
	}
}

func (e *EC2PublicKey) Verify(sig []byte, sigData []byte) (bool, error) {
	var curve elliptic.Curve

	switch e.Alg {
	case -7:
		curve = elliptic.P256()
	default:
		return false, fmt.Errorf("Error: UnSupported Algorithm")
	}

	pubkey := &ecdsa.PublicKey{
		Curve: curve,
		X:     big.NewInt(0).SetBytes(e.X),
		Y:     big.NewInt(0).SetBytes(e.Y),
	}

	ecdsaSig := &ECDSASignature{}
	_, err := asn1.Unmarshal(sig, ecdsaSig)
	if err != nil {
		return false, fmt.Errorf("Error: Invalid Signature.")
	}

	hash := crypto.SHA256.New()
	hash.Write(sigData)

	return ecdsa.Verify(pubkey, hash.Sum(nil), ecdsaSig.R, ecdsaSig.S), nil
}
