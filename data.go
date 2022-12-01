// Package for anything that handles data

package main

import (
	"fmt"
	"strings"
	"time"
)

type text string
type id uint64
type date time.Time
type state bool

func (i *id) String() string {
	return fmt.Sprintf("%d", i)
}

// Col defines column types in SQL Database
type Col interface {
	id | text | date | state
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

type RelationalDB interface {
	Insert(tableName string, tableRows Rows) error
	Delete(tableName string) error
	Retrieve(tables ...DBTable) error
	Update(planetscale, reddit DBTable, v verb) error
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

type Rows []*Row[id, text]
type DBTable *Table[Row[id, text]]

// ClearTable clears a table's row of its column values. Resets it basically.
func ClearTables(t ...*Table[Row[id, text]]) {
	for _, table := range t {
		for _, row := range table.Rows {
			if row == nil {
				continue
			}
			row.Col1 = 0
			row.Col2 = ""
			row.Col3 = ""
			row.Col4 = ""
			row.Col5 = ""
		}
	}
}
