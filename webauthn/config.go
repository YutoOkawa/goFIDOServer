package webauthn

type Config struct {
	ChallengeSize           int
	RpName                  string
	RpId                    string
	RpOrigin                string
	Timeout                 int
	RequireResidentKey      bool
	AuthenticatorAttachment string
	UserVerification        string
	Attestation             string
	CryptoParams            struct {
		Type string
		Alg  int
	}
}

var config Config

func init() {
	config.ChallengeSize = 64
	config.RpName = "FIDO_SERVER"
	config.RpId = "localhost/#/register"
	config.RpOrigin = "localhost/#/register"
	config.Timeout = 60000
	config.RequireResidentKey = false
	config.AuthenticatorAttachment = "cross-platform"
	config.UserVerification = "preferred"
	config.Attestation = "direct"
	config.CryptoParams.Type = "public-key"
	config.CryptoParams.Alg = -7
}
