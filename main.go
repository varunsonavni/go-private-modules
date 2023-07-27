package main

import (
	"log"
	"net/http"
)

func main() {
	// fmt.Println("HTTP request server listening on port 8080..")
	log.Println("HTTP request server listening on port 8080..")
	http.HandleFunc("/api", HandleAPIRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
