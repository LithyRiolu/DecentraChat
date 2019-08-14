package witai

import (
	//"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

// Conn for simple wit connect
func Conn(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	textInput, _ := vars["textInput"]

	//if err1 != nil {
	//	w.WriteHeader(400)
	//
	//}

	client := &http.Client{
		/*CheckRedirect: redirectPolicyFunc,*/
	}

	req, err := http.NewRequest("GET", "https://api.wit.ai/"+textInput, nil)
	req.Header.Add("Authorization", `Bearer GDXSJ3BC2B5JTZDZFSPNPKJWILZKLSX5`)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Err in response from wit")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//client.Do()

	w.WriteHeader(200)
	//_, _ = w.Write([]byte("hello wit ai Conn func"))
	_, _ = w.Write([]byte(body))
}
