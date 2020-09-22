package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YutoOkawa/goFIDOServer/webauthn"
)

type ErrorMessage struct {
	Message string `json:"ErrorMessage"`
}

func AssertionOptions(w http.ResponseWriter, r *http.Request) {
	var req webauthn.AuthUserRequest

	fmt.Println("-----/assertion/options-----")
	// リクエストパラメータの取得
	if err := json.Unmarshal(getReqBody(r), &req); err != nil {
		log.Fatal(err)
	}

	options, err := webauthn.AssertionOptions(req)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&ErrorMessage{Message: err.Error()})
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(options)
	}
}
