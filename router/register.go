package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YutoOkawa/goFIDOServer/webauthn"
)

type resultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func AttestationOptions(w http.ResponseWriter, r *http.Request) {
	var req webauthn.UserRequest

	fmt.Println("-----/attestation/options-----")
	// リクエストパラメータの取得
	err := json.Unmarshal(getReqBody(r), &req)
	if err != nil {
		// TODO: もっといいエラー処理
		log.Fatal(err)
	}

	// レスポンスの設定
	options, err := webauthn.AttestationOptions(req)
	if err != nil {
		log.Fatal(err)
	}

	// レスポンスパラメータの設定
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(options)
}

func AttestationResult(w http.ResponseWriter, r *http.Request) {
	var req webauthn.NavigatorCreate
	var res resultResponse

	fmt.Println("-----/attestation/result-----")
	// リクエストパラメータの取得
	err := json.Unmarshal(getReqBody(r), &req)
	if err != nil {
		log.Fatal(err)
	}

	if err = webauthn.AttestationResult(req); err != nil {
		fmt.Println(err)
	}

	res.Code = -1
	res.Message = err.Error()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
