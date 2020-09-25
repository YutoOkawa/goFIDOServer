package webauthn

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/YutoOkawa/goFIDOServer/db"
)

const (
	UserPresentFlag byte = 1 << iota
	_
	UserVerifiedFlag
	_
	_
	_
	AttestedCredentialDataFlag
	ExtenstionDataFlag
)

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
