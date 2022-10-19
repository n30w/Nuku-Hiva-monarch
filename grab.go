package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/theckman/yacspin"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// ReadAllRedditSaved reads all cached posts on the Reddit account.
// This can be used to mass refresh an entire SQL database.
func ReadAllRedditSaved(postsTable, commentsTable *Table[Row[id, text]], key *Key) {

	spinner, _ := yacspin.New(yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[43],
		Suffix:          " retrieving posts and comments",
		SuffixAutoColon: true,
		Message:         "", // Set this to the page "after" setting from struct
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	})

	var mySavedPosts []*reddit.Post
	var mySavedComments []*reddit.Comment
	var response *reddit.Response

	// Last position of for loops
	lastPos1 := 0
	lastPos2 := 0

	requests := 1000 / 100 // Reddit only caches 1000 posts
	ctx := context.Background()
	opts := &reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit:  100,
			After:  "",
			Before: "",
		},
		Sort: "new",
		Time: "all",
	}

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, err := reddit.NewClient(*key.NewKey(), reddit.WithHTTPClient(httpClient))

	if err != nil {
		log.Println("Login failed :(")
	} else {
		fmt.Println("Contacting Reddit API...")
	}

	spinner.Start()

	for i := 0; i < requests; i++ {
		mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, opts)
		if err != nil {
			log.Fatalln(err)
		}

		for _, post := range mySavedPosts {
			postsTable.Rows[lastPos1] = &Row[id, text]{
				0,
				text(post.Title),
				text(post.Permalink),
				text(post.SubredditName),
				text(post.URL),
			}
			lastPos1++
		}

		for _, comment := range mySavedComments {
			commentsTable.Rows[lastPos2] = &Row[id, text]{
				0,
				text(comment.Author),
				text(comment.Body),
				text(comment.Permalink),
				text(comment.SubredditName),
			}
			lastPos2++
		}

		spinner.Message(opts.ListOptions.After)

		opts.ListOptions.After = response.After
		time.Sleep(1 * time.Second) // Its recommend to hit Reddit with only 1 request/sec
	}

	// Populate post IDs
	for i, post := range postsTable.Rows {
		if post == nil {
			break
		}
		post.Col1 = id(lastPos1 - i)
	}

	// Populate comment IDs
	for i, comment := range commentsTable.Rows {
		if comment == nil {
			break
		}
		comment.Col1 = id(lastPos2 - i)
	}

	spinner.Stop()

	fmt.Println("Saved posts and comments retrieved")
	if err != nil {
		log.Println(err)
	}
}

// ReadRecentRedditSaved reads up to the most recent 25 saved items.
// It drops local table rows.
func ReadRecentRedditSaved(postsTable, commentsTable *Table[Row[id, text]], key *Key) {

}
