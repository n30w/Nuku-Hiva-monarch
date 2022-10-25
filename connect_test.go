package main

import (
	"os"
	"testing"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func TestNewKey(t *testing.T) {
	k := &Key{}
	got := *k.NewKey()
	want := reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	if got != want {
		t.Error("keys retrieved via .env do not match")
	}
}
