// Package for anything that handles data

package main

import (
	"fmt"
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
	[1000]*Row[id, text]
}

// Row defines a row in an SQL table
type Row[I id, T text] struct {
	Col1 I
	Col2 T
	Col3 T
	Col4 T
	Col5 T
}

// Table represents an SQL table: it has a name and rows
type Table[T Row[id, text]] struct {
	Name string
	Rows [1000]*T
}

// List prints out slice items to console
func (t *Table[T]) List() {
	for _, row := range t.Rows {
		if row == nil {
			break
		}
		fmt.Println(*row)
	}
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
