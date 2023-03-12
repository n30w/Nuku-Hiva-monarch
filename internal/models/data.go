// Package for anything that handles data
package models

import (
	"fmt"
	"strings"
	"time"
)

type text string
type id uint64
type date time.Time
type state bool

type Text text
type Id id
type Date date
type State state

func (i *Id) String() string {
	return fmt.Sprintf("%d", i)
}

func (s *State) String() string {
	return fmt.Sprintf("%T", s)
}

// col defines column types in SQL Database
type col interface {
	id | text | date | state
}

// Row defines a row in an SQL table
type Row[I col, T col] struct {
	Col1 I
	Col2 T
	Col3 T
	Col4 T
	Col5 T
}

// NewRow returns a row object given column values.
func NewRow(i int, col2, col3, col4, col5 string) *Row[id, text] {
	return &Row[id, text]{
		Col1: id(i),
		Col2: text(col2),
		Col3: text(col3),
		Col4: text(col4),
		Col5: text(col5),
	}
}

func (r *Row[I, T]) String() string {
	return fmt.Sprintf(
		"[%T, %T, %T, %T, %T]\n",
		&r.Col1, &r.Col2, &r.Col3, &r.Col4, &r.Col5,
	)
}

type DBTable *Table[Row[id, text]]
type Rows [1000]*Row[id, text]

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

// NewTable creates and returns a new table for use.
func NewTable(name string) *Table[Row[id, text]] {
	return &Table[Row[id, text]]{
		Name: name,
	}
}

// ClearTables clears a table's row of its column values. Resets it basically.
func ClearTables(tables ...DBTable) {
	for _, table := range tables {
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
