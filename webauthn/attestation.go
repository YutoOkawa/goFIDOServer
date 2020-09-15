package webauthn

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/YutoOkawa/goFIDOServer/db"
	"github.com/fxamacker/cbor"
)

type UserRequest struct {
	UserName    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type registerOptions struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
	Rp        struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"rp"`
	User struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	// AuthenticatorSelection struct {
	// 	RequireResidentKey      bool   `json:"requireResidentKey"`
	// 	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	// 	UserVerificatoin        string `json:"userVerification"`
	// } `json:"authenticatorSelection"`
	Attestation string `json:"attestation"`
	Timeout     int    `json:"timeout"`
}

type NavigatorCreate struct {
	UserName string `json:"username"`
	Create   struct {
		Id       string `json:"id"`
		RawId    string `json:"rawId"`
		Type     string `json:"type"`
		Response struct {
			AttestationObject string `json:"attestationObject"`
			ClientDataJSON    string `json:"clientDataJSON"`
		} `json:"response"`
	} `json:"create"`
}

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

func createOptions(userName string, displayName string) registerOptions {
	var options registerOptions

	options.Status = "ok"
	options.Challenge = makeRandom(config.ChallengeSize)
	options.Rp.Name = config.RpName
	options.Rp.Id = config.RpId
	options.User.Id = makeRandom(32)
	options.User.Name = userName
	options.User.DisplayName = displayName
	// options.AuthenticatorSelection.RequireResidentKey = config.RequireResidentKey
	// options.AuthenticatorSelection.AuthenticatorAttachment = config.AuthenticatorAttachment
	// options.AuthenticatorSelection.UserVerificatoin = config.UserVerification
	options.Attestation = config.Attestation
	options.Timeout = config.Timeout

	return options
}

func AttestationOptions(req UserRequest) (registerOptions, error) {
	options := createOptions(req.UserName, req.DisplayName)

	if err := db.InsertDB(options.Challenge, options.User.Id); err != nil {
		return registerOptions{}, err
	}

	return options, nil
}

func parseClientDataJSON(rawClientDataJSON string) (*ClientDataJSON, error) {
	var clientDataJSON ClientDataJSON

	clientDataJSONBin, err := base64.RawStdEncoding.DecodeString(rawClientDataJSON)
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

func verifyPackedFormat(att AttestationObject, clientDataHash []byte, authData AuthData) (bool, error) {
	alg, present := att.AttStmt["alg"].(int64)
	if !present {
		return false, fmt.Errorf("Error alg value %d\n", alg)
	}

	sig, present := att.AttStmt["sig"].([]byte)
	if !present {
		return false, fmt.Errorf("Error signature value %x\n", sig)
	}

	x5c, x509Present := att.AttStmt["x5c"].([]interface{})
	if x509Present {
		return false, fmt.Errorf("Not Implemented Error: x509 %s", x5c)
	}

	ecdaaKeyID, ecdaaKeyPresent := att.AttStmt["ecdaaKeyId"].([]byte)
	if ecdaaKeyPresent {
		return false, fmt.Errorf("Not Implemented Error: ecdaa %x", ecdaaKeyID)
	}

	return verifySelfAttestation(alg, sig, att.AuthData, clientDataHash, authData.attestedCredentialData.credentialPublicKey)
}

func verifySelfAttestation(alg int64, sig []byte, authData []byte, clientDataHash []byte, pubKey []byte) (bool, error) {
	sigData := append(authData, clientDataHash...)
	fmt.Println(sigData)

	// TODO: 公開鍵を作成する
	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return false, err
	}
	fmt.Println(publicKey)

	// TODO: 署名を検証する
	switch publicKey.(type) {
	case EC2PublicKey:
		e := publicKey.(EC2PublicKey)
		return e.Verify(sig, sigData)
	}
	return false, fmt.Errorf("failde to Verify Attestation...")
}

func AttestationResult(create NavigatorCreate) error {
	// clientDataJSONのデコード
	clientDataJSON, err := parseClientDataJSON(create.Create.Response.ClientDataJSON)
	if err != nil {
		return err
	}
	// fmt.Println(*clientDataJSON)

	// challengeの検証
	user, err := db.GetOneDB(*&clientDataJSON.Challenge)
	if err != nil {
		return err
	}
	fmt.Println("get User", user)

	// attestationObjectのデコード
	attestationObject, err := parseAttestationObject(create.Create.Response.AttestationObject)
	if err != nil {
		return err
	}
	// fmt.Println(attestationObject)

	// authenticatorDataのパース
	authData := parseAuthData(attestationObject.AuthData)
	// fmt.Println(authData)

	// TODO: attestationの検証
	clientDataHash := sha256.Sum256([]byte(create.Create.Response.ClientDataJSON))
	verify, err := verifyPackedFormat(*attestationObject, clientDataHash[:], authData)
	if err != nil {
		return err
	}

	if !verify {
		return fmt.Errorf("failed to Verify Attestation")
	}

	// TODO: 各種パラメータの検証
	// 1:originの検証
	if *&clientDataJSON.Origin != config.RpOrigin {
		return fmt.Errorf("failed to check origin!")
	}
	// 2:rpIdの検証
	rpIdHash := sha256.Sum256([]byte(config.RpId))
	if hex.EncodeToString(authData.rpIDHash) != hex.EncodeToString(rpIdHash[:]) {
		return fmt.Errorf("failed to check rpidHash")
	}

	// 3:typeの検証
	if *&clientDataJSON.Type != "webauthn.create" {
		return fmt.Errorf("failed to check type!")
	}
	// 4:flagsの検証

	// TODO: 公開鍵の作成

	// TODO: 公開鍵をデータベースに格納

	return nil
}
