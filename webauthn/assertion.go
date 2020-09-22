package webauthn

import "github.com/YutoOkawa/goFIDOServer/db"

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

func createAssertionOptions(userID string) (authOptions, error) {
	var options authOptions
	options.Challenge = makeRandom(config.ChallengeSize)
	options.RpId = config.RpId
	options.UserVerification = "require"
	options.Timeout = config.Timeout

	pubkey, err := db.GetPublicKey(userID)
	if err != nil {
		options.Status = "ng"
		options.ErrorMessage = err.Error()
		return options, err
	}

	options.AllowCredentials = make([]AllowCredential, 1)
	options.AllowCredentials[0].Id = pubkey.Keyid
	options.AllowCredentials[0].Type = "public-key"

	options.Status = "ok"
	options.ErrorMessage = ""
	return options, nil
}

func AssertionOptions(req AuthUserRequest) (authOptions, error) {
	options, err := createAssertionOptions(req.UserName)
	if err != nil {
		return options, err
	}

	// TODO: IDがユニーク制約に引っかかることを解消
	if err := db.InsertChallenge(options.Challenge, req.UserName); err != nil {
		return authOptions{}, err
	}
	return options, nil
}