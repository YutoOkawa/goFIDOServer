package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"

	"github.com/ugorji/go/codec"
)

type AttestationObject struct {
	Fmt      string  `codec:"fmt" cbor:"fmt"`
	AttStmt  AttStmt `codec:"attStmt" cbor:"attStmt"`
	AuthData []byte  `codec:"authData" cbor:"authData"`
}

type AttStmt struct {
	Sig []byte `codec:"sig" cbor:"sig"`
	// Alg []byte `codec:"alg" cbor:"alg"`
	EcdaaKeyId []byte `codec:"ecdaaKeyId" cbor:"ecdaaKeyId"`
}

type AuthData struct {
	rpIDHash               []byte
	flags                  byte
	signCount              uint32
	attestedCredentialData AttestedCredentialData
}

type AttestedCredentialData struct {
	aaguid              []byte
	credIDLen           uint16
	credID              []byte
	credentialPublicKey []byte
}

func makeRandom(i int) string {
	b := make([]byte, i)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func parseAttestationObject(rawAttestationObject string) (*AttestationObject, error) {
	var attestationObject AttestationObject
	var ch codec.CborHandle

	attestationBin, err := base64.RawURLEncoding.DecodeString(rawAttestationObject)
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoderBytes(attestationBin, &ch)
	if err := dec.Decode(&attestationObject); err != nil {
		return nil, err
	}

	return &attestationObject, nil
}

func parseAuthData(authData []byte) AuthData {
	parseAuthData := AuthData{}
	parseAuthData.rpIDHash = authData[:32]
	parseAuthData.flags = authData[32]
	signCount := authData[33:37]
	parseAuthData.signCount = binary.BigEndian.Uint32(signCount)

	parseAttestedCred := AttestedCredentialData{}
	parseAttestedCred.aaguid = authData[37:53]
	credIDLen := authData[53:55]
	parseAttestedCred.credIDLen = binary.BigEndian.Uint16(credIDLen)
	parseAttestedCred.credID = authData[55 : 55+parseAttestedCred.credIDLen]
	parseAttestedCred.credentialPublicKey = authData[55+parseAttestedCred.credIDLen:]

	parseAuthData.attestedCredentialData = parseAttestedCred

	return parseAuthData
}
