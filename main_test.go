package main

import (
	"net/http"
	"testing"
)

func TestMain(t *testing.T) {
	_, err := http.Get("http://localhost:4000/areyouawake")

	if err != nil {
		t.Errorf("%s", err)
	}

}
