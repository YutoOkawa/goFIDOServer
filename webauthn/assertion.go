package webauthn

import "github.com/YutoOkawa/goFIDOServer/db"

type AuthUserRequest struct {
	UserName string `json:"username"`
}

type authOptions struct {
	Status           string `json:"status"`
	ErrorMessage     string `json:"errorMessage"`
	Challenge        string `json:"challenge"`
	RpId             string `json:"rpId"`
	AllowCredentials []struct {
		Id   string `json:"id"`
		Type string `json:"type"`
	} `json:"allowCredentials"`
	UserVerification string `json:"userVerification"`
	Timeout          int    `json:"timeout"`
}

func createAssertionOptions() authOptions {
	var options authOptions
	options.Status = "ok"
	options.ErrorMessage = ""
	options.Challenge = makeRandom(config.ChallengeSize)
	options.RpId = config.RpId
	options.UserVerification = "require"
	options.Timeout = config.Timeout

	return options
}

func AssertionOptions(req AuthUserRequest) (authOptions, error) {
	options := createAssertionOptions()
	if err := db.InsertChallenge(options.Challenge, req.UserName); err != nil {
		return authOptions{}, err
	}
	return options, nil
}
