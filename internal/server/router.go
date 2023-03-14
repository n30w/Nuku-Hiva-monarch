package server

import (
	"fmt"
	"html/template"
	"log"
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

// homeHandler handles incoming "/" requests
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./web/template/home.tmpl.html")
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
