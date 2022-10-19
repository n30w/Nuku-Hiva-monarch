package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/theckman/yacspin"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Key represents credentials used to login to APIs
type Key struct{}

// NewKey returns a new key given environment variables
func (k *Key) NewKey() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
}

// ReadAllRedditSaved reads all cached posts on the Reddit account.
// This can be used to mass refresh an entire SQL database.
func ReadAllRedditSaved(postsTable, commentsTable *Table[Row[id, text]], key *Key) {

	ctx := context.Background()

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, err := reddit.NewClient(*key.NewKey(), reddit.WithHTTPClient(httpClient))

	if err != nil {
		log.Println("Login failed :(")
	} else {
		fmt.Println("Contacting Reddit API...")
	}

	spinner, _ := yacspin.New(yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[43],
		Suffix:          " retrieving posts and comments",
		SuffixAutoColon: true,
		Message:         "", // Set this to the page "after" setting from struct
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	})

	opts := reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit:  100,
			After:  "",
			Before: "",
		},
		Sort: "new",
		Time: "all",
	}

	// Returns for client.User.Saved method
	var mySavedPosts []*reddit.Post
	var mySavedComments []*reddit.Comment
	var response *reddit.Response

	totalRequestLimit := 1000 / 100 // Reddit only caches 1000 posts

	// Counters to keep track of current position in inner loops
	var postCounter, commentCounter id

	spinner.Start()
	for i := 0; i < totalRequestLimit; i++ {
		// Retrieved saved posts; comments
		mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, &opts)
		if err != nil {
			log.Fatalln(err)
		}

		for _, post := range mySavedPosts {
			postsTable.Rows = append(postsTable.Rows, &Row[id, text]{
				0,
				text(post.Title),
				text(post.Permalink),
				text(post.SubredditName),
				text(post.URL),
			})
			postCounter++
		}

		for _, comment := range mySavedComments {
			commentsTable.Rows = append(commentsTable.Rows, &Row[id, text]{
				0,
				text(comment.Author),
				text(comment.Body),
				text(comment.Permalink),
				text(comment.SubredditName),
			})
			commentCounter++
		}

		spinner.Message(opts.ListOptions.After)

		// Update ListOptions.After
		opts.ListOptions.After = response.After
		time.Sleep(1 * time.Second) // Its recommend to hit Reddit with only 1 request/sec
	}

	lenPosts := len(postsTable.Rows)
	lenComments := len(commentsTable.Rows)

	// Populate post IDs
	for i, post := range postsTable.Rows {
		post.Col1 = id(lenPosts - i)
	}

	// Populate comment IDs
	for i, comment := range commentsTable.Rows {
		comment.Col1 = id(lenComments - i)
	}

	spinner.Stop()
	fmt.Println("Saved posts and comments retrieved")
	if err != nil {
		log.Println(err)
	}
}

// ReadRecentRedditSaved reads up to the most recent 25 saved items.
// It drops local table rows.
func ReadRecentRedditSaved() {

}
