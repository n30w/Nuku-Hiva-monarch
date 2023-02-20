package reddit

import (
	"testing"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
)

// func TestPopulateIDs(t *testing.T) {

// 	var mockTable DBTable
// 	mockTable.Name = "mockTable"

// }

func BenchmarkGrabSaved(b *testing.B) {
	postsTable := models.NewTable("posts")
	commentsTable := models.NewTable("comments")
	key := &credentials.RedditKey{}

	for i := 0; i < b.N; i++ {
		Saved(postsTable, commentsTable, key)
	}

}
