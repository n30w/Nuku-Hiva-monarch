// Package for anything that handles data

package main

import (
	"fmt"
	"strings"
)

type text string
type id uint64

func (i *id) String() string {
	return fmt.Sprintf("%d", i)
}

// Col defines column types in SQL Database
type Col interface {
	id | text
}

// Row defines a row in an SQL table
type Row[I Col, T Col] struct {
	Col1 I
	Col2 T
	Col3 T
	Col4 T
	Col5 T
}

func (r *Row[I, T]) String() string {
	return fmt.Sprintf(
		"[%T, %T, %T, %T, %T]",
		r.Col1, r.Col2, r.Col3, r.Col4, r.Col5,
	)
}

// Table represents an SQL table: it has a name and rows
type Table[T Row[id, text]] struct {
	Name string
	Rows [1000]*Row[id, text]
}

func (t *Table[Row]) String() string {
	var sb strings.Builder
	for _, row := range t.Rows {
		if row == nil {
			break
		}
		sb.WriteString(row.String() + "\n")
	}
	return sb.String()
}
