package main

import (
	"log"
	"net/http"

	"github.com/YutoOkawa/goFIDOServer/router"
)

func main() {
	r := router.New()
	r.SetRouter()
	err := http.ListenAndServeTLS(":8080", "ssl/myself.crt", "ssl/myself.key", r.Router)
	if err != nil {
		log.Fatal(err)
	}
}
