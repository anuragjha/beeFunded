package main

import (
	"./uri_routing"
	"log"
	"net/http"
	"os"
)

func main() {

	router := uri_routing.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}

}
