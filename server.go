package main

import (
	"log"
	"net/http"
)

// Initialize initializes a web server
func Initialize() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", update)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Update Planetscale Database"))
}
