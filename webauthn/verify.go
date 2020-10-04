package webauthn

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/YutoOkawa/goFIDOServer/db"
)

type AuthenticatorFlag byte

const (
	UserPresentFlag AuthenticatorFlag = 1 << iota
	_
	UserVerifiedFlag
	_
	_
	_
	AttestedCredentialDataFlag
	ExtenstionDataFlag
)

func (flag AuthenticatorFlag) VerifyUserPresent() bool {
	return (flag & UserPresentFlag) == UserPresentFlag
}

func (flag AuthenticatorFlag) VerifyUserVerified() bool {
	return (flag & UserVerifiedFlag) == UserVerifiedFlag
}

func (flag AuthenticatorFlag) HasAttestedCredentialData() bool {
	return (flag & AttestedCredentialDataFlag) == AttestedCredentialDataFlag
}

func (flag AuthenticatorFlag) HasExtensionData() bool {
	return (flag & ExtenstionDataFlag) == ExtenstionDataFlag
}

func verifyChallenge(challenge string) error {
	user, err := db.GetChallenge(challenge)
	if err != nil {
		return fmt.Errorf("failed to verify challenge%s %s", user.Challenge, err.Error())
	}
	return nil
}

func verifyParameters(clientDataJSON ClientDataJSON, authData AuthData, webauthnType string) error {
	// TODO: 各種パラメータの検証
	// 1:originの検証
	if clientDataJSON.Origin != config.RpOrigin {
		return fmt.Errorf("failed to check origin!")
	}

	// 2:rpIdの検証
	rpIdHash := sha256.Sum256([]byte(config.RpId))
	if hex.EncodeToString(authData.rpIDHash) != hex.EncodeToString(rpIdHash[:]) {
		return fmt.Errorf("failed to check rpidHash")
	}

	// 3:typeの検証
	if *&clientDataJSON.Type != webauthnType {
		return fmt.Errorf("failed to check type!")
	}

	// 4:flagsの検証
	return nil
}
