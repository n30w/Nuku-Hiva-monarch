package main

import (
	"net/http"
	"testing"
)

func TestMain(t *testing.T) {
	got, err := http.Get("http://localhost:4000/areyouawake")
	want := 200

	if err != nil {
		t.Errorf("%s", err)
	}

	if got.StatusCode != want {
		t.Errorf("server is unreachable")
	}

}
