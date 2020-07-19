package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type userRequest struct {
	UserName    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type registerOptions struct {
	Status    string `json:"status"`
	Challenge string `json:"challenge"`
	Rp        struct {
		Name string `json:"name"`
	} `json:"rp"`
	User struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
	} `json:"user"`
	AuthenticatorSelection struct {
		RequireResidentKey      bool   `json:"requireResidentKey"`
		AuthenticatorAttachment string `json:"authenticatorAttachment"`
		UserVerificatoin        string `json:"userVerification"`
	} `json:"authenticatorSelection"`
	Attestation string `json:"attestation"`
	Timeout     int    `json:"timeout"`
}

func attestationOptions(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	var options registerOptions

	fmt.Println("-----/attestation/options-----")
	// リクエストパラメータの取得
	body := r.Body
	defer body.Close()
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	json.Unmarshal(buf.Bytes(), &req)
	fmt.Println(req)

	// レスポンスの設定
	options.Status = "ok"
	options.Challenge = makeRandom(config.ChallengeSize)
	options.Rp.Name = config.RpName
	options.User.Id = makeRandom(32)
	options.User.Name = req.UserName
	options.User.DisplayName = req.DisplayName
	options.AuthenticatorSelection.RequireResidentKey = config.RequireResidentKey
	options.AuthenticatorSelection.AuthenticatorAttachment = config.AuthenticatorAttachment
	options.AuthenticatorSelection.UserVerificatoin = config.UserVerification
	options.Attestation = config.Attestation
	options.Timeout = config.Timeout

	// レスポンスパラメータの設定
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(options)
}

func attestationResult() {

}
