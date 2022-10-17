package main

// Package for anything that handles data

import (
	"database/sql"
	"log"
	"strconv"
)

type text string
type id uint64

type Table[T Row[id, text]] struct {
	Name string
	Rows []*T
}

// Different types of rows of different tables
type Rows interface {
	[]*Row[id, text]
}

// A column can either be an int or a string. Ints are for IDs.
type Col interface {
	id | text
}

type Row[I id, T text] struct {
	Col1 I
	Col2 T
	Col3 T
	Col4 T
	Col5 T
}

func RetrieveBy[T Col](db *sql.DB, table *Table[Row[id, text]], retrieval, want *T) T {
	q := "SELECT " + ToString(retrieval) + " FROM " + table.Name + " WHERE " + ToString(want)

	err := db.QueryRow(q).Scan(&retrieval)
	if err == sql.ErrNoRows {
		log.Fatal("no rows returned")
	} else if err != nil {
		log.Fatal(err)
	}
	return *retrieval
}

func ToString(t any) string {
	switch v := t.(type) {
	case string:
		return v
	case id:
		return strconv.Itoa(int(v))
	default:
		return ""
	}
}
