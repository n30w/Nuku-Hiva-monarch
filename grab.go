package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/theckman/yacspin"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func ReadRedditData() (*Table[Row[id, text]], *Table[Row[id, text]]) {

	// posts := make([]*Post, 0)
	// comments := make([]*Comment, 0)

	postsTable := &Table[Row[id, text]]{}
	commentsTable := &Table[Row[id, text]]{}

	postsTable.Name = "posts"
	commentsTable.Name = "comments"

	// Retrieves credentials from txt file. Please figure out a
	// way to encrypt this info. Thanks.
	credentials := func() reddit.Credentials {

		line := 0
		c := [4]string{}

		f, err := os.Open("credentials.txt")

		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		defer f.Close()

		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			if line != 0 {
				switch line {
				case 1:
					c[2] = scanner.Text()
				case 2:
					c[3] = scanner.Text()
				case 3:
					c[0] = scanner.Text()
				case 4:
					c[1] = scanner.Text()
				}
			}
			line++
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		return reddit.Credentials{
			ID:       c[0],
			Secret:   c[1],
			Username: c[2],
			Password: c[3],
		}
	}()

	fmt.Println("Credentials retrieved from file")

	var ctx = context.Background()

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, err := reddit.NewClient(credentials, reddit.WithHTTPClient(httpClient))

	if err != nil {
		fmt.Println("Login failed :(")
		return nil, nil
	} else {
		fmt.Println("Contacting Reddit API...")
	}

	spinner, _ := yacspin.New(yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[43],
		Suffix:          " retrieving posts",
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
			log.Fatal(err)
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
		fmt.Println(err)
		return nil, nil
	}

	return postsTable, commentsTable
}

// TODO: GoRoutines for pulling from Reddit
// TODO: Figure out cryptography solution

/*

//Print out saved posts
	{
		title := color.New(color.FgHiGreen)
		link := color.New(color.FgCyan)
		subreddit := color.New(color.FgHiRed)
		fmt.Println("")
		for _, post := range savedPosts {
			title.Printf("\n# %s", post.Title)
			fmt.Print(" in ")
			subreddit.Printf("%s\n", post.SubredditName)
			link.Printf("- %s\n- %s\n", post.Permalink, post.URL)
		}

		fmt.Println("===========================")

		author := color.New(color.FgHiGreen)

		for _, comment := range savedComments {
			author.Printf("\n@%s", comment.Author)
			fmt.Print(" in ")
			subreddit.Printf("%s\n", comment.SubredditName)
			fmt.Printf("%s\n", comment.Body)
			fmt.Printf("\n%s\n", comment.Permalink)
		}
	}



*/
