package reddit

import "testing"

// func TestPopulateIDs(t *testing.T) {

// 	var mockTable DBTable
// 	mockTable.Name = "mockTable"

// }

func BenchmarkGrabSaved(b *testing.B) {
	postsTable := &Table[Row[id, text]]{}
	commentsTable := &Table[Row[id, text]]{}
	key := &Key{}

	for i := 0; i < b.N; i++ {
		GrabSaved(postsTable, commentsTable, key)
	}

}
