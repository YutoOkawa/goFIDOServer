package router

import (
	"bytes"
	"io"
	"net/http"
)

func getReqBody(req *http.Request) []byte {
	body := req.Body
	defer body.Close()
	buf := new(bytes.Buffer)
	io.Copy(buf, body)
	return buf.Bytes()
}
