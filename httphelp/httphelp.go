package httphelp

import (
	"io/ioutil"
	"log"
	"net/http"
)

func ReadHttpRequestBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error in RecvPeerAlive : ", err)
	}
	defer r.Body.Close()
	return body
}
