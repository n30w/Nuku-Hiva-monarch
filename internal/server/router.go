package server

import (
	"net/http"
)

func scanAndDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

// updateHandler is an http handler that handles requests to update the database.
func updateHandler(f func() error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write([]byte("database successfully updated.")) //nolint

	}
}
