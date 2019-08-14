package main

import (
	"./routing"
	"log"
	"net/http"
	"os"
)

func main() {

	router := routing.NewRouter()

	if len(os.Args) > 1 {
		log.Println("started listening on port : ", os.Args[1])
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))

	} else {
		log.Println("started listening on port : 8080")
		log.Fatal(http.ListenAndServe(":8080", router))
	}

}
