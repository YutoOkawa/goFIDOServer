package webauthn

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/YutoOkawa/goFIDOServer/db"
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

	if err := db.InsertChallenge(options.Challenge, options.User.Id); err != nil {
		return registerOptions{}, err
	}

	return options, nil
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

	// 公開鍵を作成する
	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return false, err
	}

	// 署名を検証する
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

	// challengeの検証
	if err := verifyChallenge(*&clientDataJSON.Challenge); err != nil {
		return err
	}

	// attestationObjectのデコード
	attestationObject, err := parseAttestationObject(create.Create.Response.AttestationObject)
	if err != nil {
		return err
	}

	// authenticatorDataのパース
	authData := parseAuthData(attestationObject.AuthData)

	// Attestationの検証
	clientData, err := base64.RawURLEncoding.DecodeString(create.Create.Response.ClientDataJSON)
	if err != nil {
		return fmt.Errorf("failed to decode clientDataJSON")
	}
	clientDataHash := sha256.Sum256(clientData)
	verify, err := verifyPackedFormat(*attestationObject, clientDataHash[:], authData)
	if err != nil {
		return err
	}
	if !verify {
		return fmt.Errorf("failed to Verify Attestation")
	}

	// 各種パラメータの検証
	if err := verifyParameters(*clientDataJSON, authData, "webauthn.create"); err != nil {
		return err
	}

	// 公開鍵の作成
	publicKey, err := parsePublicKey(authData.attestedCredentialData.credentialPublicKey)
	if err != nil {
		return err
	}

	// challengeからユーザIDを取得
	user, err := db.GetChallenge(*&clientDataJSON.Challenge)
	if err != nil {
		return err
	}
	userId := user.Userid

	// 公開鍵をデータベースに格納
	switch publicKey.(type) {
	case EC2PublicKey:
		e := publicKey.(EC2PublicKey)
		ec2Bytes, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if err := db.InsertPublicKey(create.Create.Id, userId, create.UserName, ec2Bytes); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Not Implemented Error")
	}

	// 認証回数を格納
	if err := db.InsertUserData(userId, create.UserName, authData.signCount); err != nil {
		if err := db.UpdateSignCount(userId, authData.signCount); err != nil {
			return err
		}
	}

	return nil
}
