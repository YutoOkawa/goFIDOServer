package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YutoOkawa/goFIDOServer/webauthn"
)

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

	fmt.Println("-----/attestation/result-----")
	// リクエストパラメータの取得
	err := json.Unmarshal(getReqBody(r), &req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)

	err = webauthn.AttestationResult(req)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
