package models

import (
	"fmt"
	"testing"
)

func TestClearTables(t *testing.T) {

	table := &Table[Row[Id, Text]]{
		Name: "test table",
	}

	comparison := [1000]int{}

	// Fill with dummy data
	for i := 0; i < 1000; i++ {
		table.Rows[i] = &Row[Id, Text]{
			Col1: Id(i),
			Col2: Text(fmt.Sprintf("%d", i)),
			Col3: Text(fmt.Sprintf("%d", i)),
			Col4: Text(fmt.Sprintf("%d", i)),
			Col5: Text(fmt.Sprintf("%d", i)),
		}
		comparison[i] = i
	}

	ClearTables(table)

	// Compare data
	for i := 1; i < 1000; i++ {
		if table.Rows[i].Col1 == Id(comparison[i]) {
			t.Errorf("got occupied table, data remains at index %d", i)
		}
	}
}
