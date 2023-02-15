package models

import (
	"fmt"
	"testing"
)

func TestClearTables(t *testing.T) {

	table := NewTable("test table")

	comparison := [1000]int{}

	// Fill with dummy data
	for i := 0; i < 1000; i++ {
		table.Rows[i] = NewRow(
			i,
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%d", i),
		)
		comparison[i] = i
	}

	ClearTables(table)

	// Compare data
	for i := 1; i < 1000; i++ {
		if table.Rows[i].Col1 == id(comparison[i]) {
			t.Errorf("got occupied table, data remains at index %d", i)
		}
	}
}
