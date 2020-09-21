package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YutoOkawa/goFIDOServer/webauthn"
)

func AssertionOptions(w http.ResponseWriter, r *http.Request) {
	var req webauthn.AuthUserRequest

	fmt.Println("-----/assertion/options-----")
	// リクエストパラメータの取得
	if err := json.Unmarshal(getReqBody(r), &req); err != nil {
		log.Fatal(err)
	}

	options, err := webauthn.AssertionOptions(req)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(options)
}
