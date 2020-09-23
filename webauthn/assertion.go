package webauthn

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/YutoOkawa/goFIDOServer/db"
)

type AuthUserRequest struct {
	UserName string `json:"username"`
}

type AllowCredential struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type authOptions struct {
	Status           string            `json:"status"`
	ErrorMessage     string            `json:"errorMessage"`
	Challenge        string            `json:"challenge"`
	RpId             string            `json:"rpId"`
	AllowCredentials []AllowCredential `json:"allowCredentials"`
	UserVerification string            `json:"userVerification"`
	Timeout          int               `json:"timeout"`
}

type NavigatorGet struct {
	UserName string `json:"username"`
	Get      struct {
		Id       string `json:"id"`
		RawId    string `json:"rawId"`
		Type     string `json:"type"`
		Response struct {
			AuthenticatorData string `json:"authenticatorData"`
			ClientDataJSON    string `json:"clientDataJSON"`
			Signature         string `json:"signature"`
			UserHandle        string `json:"userHandle"`
		} `json:"response"`
	} `json:"get"`
}

func createAssertionOptions(userID string) (authOptions, string, error) {
	var options authOptions
	options.Challenge = makeRandom(config.ChallengeSize)
	options.RpId = config.RpId
	options.UserVerification = "required"
	options.Timeout = config.Timeout

	pubkey, err := db.GetPublicKey(userID)
	if err != nil {
		options.Status = "ng"
		options.ErrorMessage = err.Error()
		return options, "", err
	}

	options.AllowCredentials = make([]AllowCredential, 1)
	options.AllowCredentials[0].Id = pubkey.Keyid
	options.AllowCredentials[0].Type = "public-key"

	options.Status = "ok"
	options.ErrorMessage = ""
	return options, pubkey.Userid, nil
}

func AssertionOptions(req AuthUserRequest) (authOptions, error) {
	options, userId, err := createAssertionOptions(req.UserName)
	if err != nil {
		return options, err
	}

	if err := db.InsertChallenge(options.Challenge, userId); err != nil {
		return authOptions{}, err
	}
	return options, nil
}

func DeleteChallenge(challenge string, retErr error) error {
	if err := db.DeleteChallenge(challenge); err != nil {
		return err
	}

	return retErr
}

func AssertionResult(get NavigatorGet) error {
	// clientDataJSONのデコード
	clientDataJSON, err := parseClientDataJSON(get.Get.Response.ClientDataJSON)
	if err != nil {
		return err
	}

	// challengeの検証
	if err := verifyChallenge(*&clientDataJSON.Challenge); err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}

	// authenticatorDataのデコード
	authDataBin, err := base64.RawURLEncoding.DecodeString(get.Get.Response.AuthenticatorData)
	if err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}
	authData := parseAuthData(authDataBin, false)

	// 各種パラメータの検証
	if err := verifyParameters(*clientDataJSON, authData, "webauthn.get"); err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}

	// 公開鍵の取得
	pubkeyData, err := db.GetPublicKey(get.UserName)
	if err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}
	var pubkey interface{}
	if err := json.Unmarshal(pubkeyData.Publickey, &pubkey); err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}

	// 署名検証
	clientData, err := base64.RawStdEncoding.DecodeString(get.Get.Response.ClientDataJSON)
	clientDataHash := sha256.Sum256(clientData)
	sigData := append(authDataBin, clientDataHash[:]...)
	signature, err := base64.RawURLEncoding.DecodeString(get.Get.Response.Signature)
	if err != nil {
		return DeleteChallenge(clientDataJSON.Challenge, err)
	}

	switch pubkey.(type) {
	case EC2PublicKey:
		e := pubkey.(EC2PublicKey)
		check, err := e.Verify(signature, sigData)
		if !check && err != nil {
			return DeleteChallenge(clientDataJSON.Challenge, err)
		}
	}

	// TODO: 認証回数の更新

	return DeleteChallenge(clientDataJSON.Challenge, nil)
}
