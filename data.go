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

func (s *state) String() string {
	return fmt.Sprintf("%T", s)
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
		"[%T, %T, %T, %T, %T]\n",
		&r.Col1, &r.Col2, &r.Col3, &r.Col4, &r.Col5,
	)
}

type DBTable *Table[Row[id, text]]
type Rows [10000]*Row[id, text]

type RelationalDB interface {
	Insert(tableName string, tableRows Rows) error
	Delete(tableName string) error
	Retrieve(amount amount, tables ...DBTable) error
	Update(planetscale, reddit DBTable, v verb) error
}

// Table represents an SQL table: it has a name and rows
type Table[T Row[id, text]] struct {
	Name string
	Rows Rows
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

// ClearTables clears a table's row of its column values. Resets it basically.
func ClearTables(t ...DBTable) { // TODO go routine optimization can occur here
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

// CreateTable instantiates and returns a pointer to a populated table with zero values.
func CreateTable(name string) DBTable {
	return createTable(name)
}

func createTable(name string) DBTable {
	table := &Table[Row[id, text]]{Name: name}

	for i := 0; i < len(table.Rows); i++ {
		table.Rows[i] = &Row[id, text]{
			0,
			"",
			"",
			"",
			"",
		}
	}

	return table
}
