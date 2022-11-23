package main

import (
	"fmt"
	"testing"
)

func TestClearTable(t *testing.T) {

	table := &Table[Row[id, text]]{
		Name: "test table",
	}

	comparison := [1000]int{}

	// Fill with dummy data
	for i := 0; i < 1000; i++ {
		table.Rows[i] = &Row[id, text]{
			Col1: id(i),
			Col2: text(fmt.Sprintf("%d", i)),
			Col3: text(fmt.Sprintf("%d", i)),
			Col4: text(fmt.Sprintf("%d", i)),
			Col5: text(fmt.Sprintf("%d", i)),
		}
		comparison[i] = i
		fmt.Println(table.Rows[i].Col1)
	}

	ClearTable(table)

	// Compare data
	for i := 1; i < 1000; i++ {
		if table.Rows[i].Col1 == id(comparison[i]) {
			t.Errorf("table not cleared, some data still remains at index %d", i)
		}
	}
}
