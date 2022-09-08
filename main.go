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
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func main() {

	/*
		Retrieves credentials from text file. This is only kept locally.
		Also, its pretty bad to keep this kind of information in plain text.
		Don't do that.
	*/
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

	var ctx = context.Background()

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, _ := reddit.NewClient(credentials, reddit.WithHTTPClient(httpClient))

	saved, err := func() ([]reddit.Post, error) {
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
		// var mySavedComments []*reddit.Comment
		var response *reddit.Response
		var err error
		// All saved posts
		allSaved := []reddit.Post{}

		// Reddit's API total request limit for saved posts
		totalRequestLimit := 1000 / 100 // The 100 is for Limit

		for i := 0; i < totalRequestLimit; i++ {
			// Retrieved saved posts; comments
			mySavedPosts, _, response, err = client.User.Saved(ctx, &opts)
			if err != nil {
				return nil, err
			}

			fmt.Println("Save pull in progress...")
			for _, post := range mySavedPosts {
				allSaved = append(allSaved, *post)
			}

			// Update ListOptions.After
			opts.ListOptions.After = response.After
			time.Sleep(1 * time.Second) // Its recommend to only hit Reddit with 1 request/sec
		}

		return allSaved, err
	}()
	if err != nil {
		return
	}

	// Print out saved posts
	{
		title := color.New(color.BgHiYellow)
		link := color.New(color.FgCyan)

		for _, post := range saved {
			title.Printf("%s", post.Title)
			fmt.Print(" | ")
			link.Printf("%s\n", post.URL)
		}
		fmt.Println("============================")
		// for _, post := range mySavedComments {
		// 	author.Printf("%s ", post.Author)
		// 	fmt.Print("in")
		// 	subredditName.Printf(" %s ", post.SubredditName)
		// 	fmt.Print("@")
		// 	commentLink.Printf(" %s\n", post.PostPermalink)
		// 	fmt.Printf("%s\n\n", post.Body)
		// }
	}

}

// TODO send custom JSON request to reddit to display more than 100 saved posts
// https://old.reddit.com/r/redditdev/comments/d7egb/how_to_get_more_json_results_i_get_only_30/
