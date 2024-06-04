package main

import (
	"log"
	"net/http"
) 

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: "127.0.0.1:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}