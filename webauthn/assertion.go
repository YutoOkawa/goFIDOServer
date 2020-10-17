package webauthn

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/YutoOkawa/goFIDOServer/db"
)

type AuthUserRequest struct {
	UserName string `json:"username"`
}

type AllowCredential struct {
	Id        string   `json:"id"`
	Type      string   `json:"type"`
	Transport []string `json:"transports"`
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
	options.AllowCredentials[0].Transport = make([]string, 4)
	options.AllowCredentials[0].Transport = []string{"usb", "internal"}

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

func AssertionResult(get NavigatorGet) error {
	// clientDataJSONのデコード
	clientDataJSON, err := parseClientDataJSON(get.Get.Response.ClientDataJSON)
	if err != nil {
		return err
	}

	// challengeの検証
	if err := verifyChallenge(*&clientDataJSON.Challenge); err != nil {
		return err
	}

	// authenticatorDataのデコード
	authDataBin, err := base64.RawURLEncoding.DecodeString(get.Get.Response.AuthenticatorData)
	if err != nil {
		return err
	}
	authData := parseAuthData(authDataBin)

	// 各種パラメータの検証
	if err := verifyParameters(*clientDataJSON, authData, "webauthn.get"); err != nil {
		return err
	}

	// 公開鍵の取得
	pubkeyData, err := db.GetPublicKey(get.UserName)
	if err != nil {
		return err
	}
	var pubkey interface{}
	if err := json.Unmarshal(pubkeyData.Publickey, &pubkey); err != nil {
		return err
	}

	// 署名検証
	clientData, err := base64.RawStdEncoding.DecodeString(get.Get.Response.ClientDataJSON)
	clientDataHash := sha256.Sum256(clientData)
	sigData := append(authDataBin, clientDataHash[:]...)
	signature, err := base64.RawURLEncoding.DecodeString(get.Get.Response.Signature)
	if err != nil {
		return err
	}

	switch pubkey.(type) {
	case EC2PublicKey:
		e := pubkey.(EC2PublicKey)
		check, err := e.Verify(signature, sigData)
		if !check && err != nil {
			return err
		}
	}

	// challengeからユーザIDを取得
	user, err := db.GetChallenge(*&clientDataJSON.Challenge)
	if err != nil {
		return err
	}
	userId := user.Userid

	// 認証回数の検証
	userData, err := db.GetUserData(userId)
	if err != nil {
		return err
	}
	if userData.Signcount > authData.signCount {
		return fmt.Errorf("failed to verify SignCount")
	}
	// 認証回数の更新
	if err := db.UpdateSignCount(userId, authData.signCount); err != nil {
		return err
	}

	return nil
}
