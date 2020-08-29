package webauthn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func AttestationResult(w http.ResponseWriter, r *http.Request) {
	var req serverRequest

	fmt.Println("-----/attestation/result-----")
	// リクエストパラメータの取得
	err := json.Unmarshal(getReqBody(r), &req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)

	clientDataJSON, err := parseClientDataJSON(req.Create.Response.ClientDataJSON)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*clientDataJSON)

	attestationObject, err := parseAttestationObject(req.Create.Response.AttestationObject)
	if err != nil {
		log.Fatal(err)
	}

	authData := parseAuthData(attestationObject.AuthData)
	fmt.Println(authData)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(clientDataJSON)
}
