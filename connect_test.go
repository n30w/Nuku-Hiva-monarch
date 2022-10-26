package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Fatal(message)
}

type MockPsDB struct{}

func (m *MockPsDB) insertToSQL(_, _ []*Row[id, text]) error {
	return nil
}

func (m *MockPsDB) RetrieveSQL(_ ...*Table[Row[id, text]]) error {
	return nil
}

func (m *MockPsDB) UpdateSQL(_, _ *Table[Row[id, text]], _ verb) error {
	return nil
}

func TestNewKey(t *testing.T) {
	k := &Key{}
	got := *k.NewKey()
	want := reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	assertEqual(t, got, want, "keys retrieved via .env do not match")

}

func TestSQL(t *testing.T) {
	var db *sql.DB
	db, err := sql.Open("mysql", os.Getenv("DEV"))
	if err != nil {
		panic(Warn.Sprint(err))
	}

	if err := db.Ping(); err != nil {
		panic(Warn.Sprint(err))
	}

	psdb := &PlanetscaleDB{db}

	// Test tables
	redditPosts := &Table[Row[id, text]]{Name: "posts"}
	dBPosts := &Table[Row[id, text]]{Name: "posts"}

	t.Run("get last ID", func(t *testing.T) {
		var got id = psdb.getLastID("posts")
		var want id = 500

		assertEqual(t, got, want, "")
	})

	t.Run()
}
