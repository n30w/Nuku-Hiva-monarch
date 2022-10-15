package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/theckman/yacspin"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func Grab() {

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
		return
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

	savedPosts, savedComments, err := func() ([]reddit.Post, []reddit.Comment, error) {
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
		var err error

		// Function return values
		allSavedPosts := []reddit.Post{}
		allSavedComments := []reddit.Comment{}
		totalRequestLimit := 1000 / 100 // Reddit only caches 1000 posts

		spinner.Start()
		for i := 0; i < totalRequestLimit; i++ {
			// Retrieved saved posts; comments
			mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, &opts)
			if err != nil {
				return nil, nil, err
			}

			for _, post := range mySavedPosts {
				allSavedPosts = append(allSavedPosts, *post)
			}
			for _, comment := range mySavedComments {
				allSavedComments = append(allSavedComments, *comment)
			}

			spinner.Message(opts.ListOptions.After)

			// Update ListOptions.After
			opts.ListOptions.After = response.After
			time.Sleep(1 * time.Second) // Its recommend to only hit Reddit with 1 request/sec

		}
		spinner.Stop()
		return allSavedPosts, allSavedComments, err
	}()
	if err != nil {
		return
	}
	// Print out saved posts
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

}

// TODO: GoRoutines for pulling from Reddit
// TODO: Figure out database solution
// TODO: Figure out cryptography solution
