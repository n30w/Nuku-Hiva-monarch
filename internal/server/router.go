package server

import (
	"fmt"
	"net/http"
)

// serverOperation is what the httpHandler calls when a request is received.
type serverOperation func() error

func scanDeleteHandler(f serverOperation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			w.Write([]byte(fmt.Sprint(err))) //nolint
			fmt.Println(err)
		} else {
			w.Write([]byte("duplicate entries scanned and deleted")) //nolint
		}
	}
}

// updateHandler is an http handler that handles requests to update the database.
func updateHandler(f serverOperation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write([]byte("database successfully updated")) //nolint
	}
}

func populateHandler(f serverOperation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write([]byte("database successfully populated")) //nolint
	}
}
