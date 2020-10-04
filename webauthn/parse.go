package webauthn

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"github.com/fxamacker/cbor"
)

type ClientDataJSON struct {
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Origin    string `json:"origin"`
}

type AttestationObject struct {
	Fmt      string                 `json:"fmt"`
	AttStmt  map[string]interface{} `json:"attStmt,omitempty"`
	AuthData []byte                 `json:"authData"`
}

type AuthData struct {
	rpIDHash               []byte
	flags                  AuthenticatorFlag
	signCount              uint32
	attestedCredentialData AttestedCredentialData
}

type AttestedCredentialData struct {
	aaguid              []byte
	credIDLen           uint16
	credID              []byte
	credentialPublicKey []byte
}

func parseClientDataJSON(rawClientDataJSON string) (*ClientDataJSON, error) {
	var clientDataJSON ClientDataJSON

	clientDataJSONBin, err := base64.RawURLEncoding.DecodeString(rawClientDataJSON)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(clientDataJSONBin, &clientDataJSON)
	if err != nil {
		return nil, err
	}

	return &clientDataJSON, nil
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

func parseAuthData(authData []byte, isKey bool) AuthData {
	parseAuthData := AuthData{}
	parseAuthData.rpIDHash = authData[:32]
	parseAuthData.flags = AuthenticatorFlag(authData[32])
	signCount := authData[33:37]
	parseAuthData.signCount = binary.BigEndian.Uint32(signCount)

	if isKey {
		parseAttestedCred := AttestedCredentialData{}
		parseAttestedCred.aaguid = authData[37:53]
		credIDLen := authData[53:55]
		parseAttestedCred.credIDLen = binary.BigEndian.Uint16(credIDLen)
		parseAttestedCred.credID = authData[55 : 55+parseAttestedCred.credIDLen]
		parseAttestedCred.credentialPublicKey = authData[55+parseAttestedCred.credIDLen:]
		parseAuthData.attestedCredentialData = parseAttestedCred
	}

	return parseAuthData
}
