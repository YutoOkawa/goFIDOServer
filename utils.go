package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	"net/http"

	"github.com/fxamacker/cbor"
)

type AttestationObject struct {
	Fmt      string                 `json:"fmt"`
	AttStmt  map[string]interface{} `json:"attStmt,omitempty"`
	AuthData []byte                 `json:"authData"`
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

func getReqBody(req *http.Request) []byte {
	body := req.Body
	defer body.Close()
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	return buf.Bytes()
}

func parseAttestationObject(rawAttestationObject string) (*AttestationObject, error) {
	var attestationObject AttestationObject

	attestationBin, err := base64.RawURLEncoding.DecodeString(rawAttestationObject)
	if err != nil {
		return nil, err
	}

	err = cbor.Unmarshal(attestationBin, &attestationObject)
	if err != nil {
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
