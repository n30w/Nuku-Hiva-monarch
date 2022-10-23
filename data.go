// Package for anything that handles data

package main

import (
	"fmt"
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
	Rows [1000]*Row[id, text]
}

// List prints out slice items to console
func (t *Table[T]) List() {
	for _, row := range t.Rows {
		if row == nil {
			fmt.Println("No rows in this table")
			break
		}
		fmt.Println(*row)
	}
}

// Flush clears all the data in each respective table
func (t *Table[T]) Flush() {
	for _, row := range t.Rows {
		*row = Row[id, text]{}
	}
}
