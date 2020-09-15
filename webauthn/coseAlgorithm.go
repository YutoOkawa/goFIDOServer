package webauthn

import (
	"fmt"

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

func parsePublicKey(pubkeyBytes []byte) (interface{}, error) {
	var pubKey PublicKey
	err := cbor.Unmarshal(pubkeyBytes, &pubKey)
	if err != nil {
		return nil, err
	}

	fmt.Println("PublicKey", pubKey)

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

	return false, fmt.Errorf("Not Implemented Error")
}
