package main

import (
	"log"
	"net/http"

	"bookapi/routes" // import the correct module path
)

func main() {
	log.Println("Starting server on http://127.0.0.1:8080...")
	if err := http.ListenAndServe("127.0.0.1:8080", routes.Routes()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
