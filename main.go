package main

import (
	"log"
	"net/http"

	"github.com/varunsonavni/go-private-modules/src/helm"
)

func main() {
	// fmt.Println("HTTP request server listening on port 8080..")
	log.Println("HTTP request server listening on port 8080..")
	http.HandleFunc("/api", helm.HandleAPIRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
