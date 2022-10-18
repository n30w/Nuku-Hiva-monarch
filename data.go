package main

// Package for anything that handles data

import (
	"strconv"
)

type text string
type id uint64

// Col defines column types in SQL Database
type Col interface {
	id | text
}

// Rows consist of IDs and text
type Rows interface {
	[]*Row[id, text]
}

// Row defines a row in an SQL table
type Row[I id, T text] struct {
	Col1 I
	Col2 T
	Col3 T
	Col4 T
	Col5 T
}

type Table[T Row[id, text]] struct {
	Name string
	Rows []*T
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
