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

	fmt.Println("Credentials retrieved!")

	var ctx = context.Background()

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, _ := reddit.NewClient(credentials, reddit.WithHTTPClient(httpClient))

	fmt.Println("Contacting Reddit API...")

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

		// Reddit's API total request limit for saved posts
		totalRequestLimit := 1000 / 100 // The 100 is for Limit

		// Loading dots
		loadingDots := [3]string{".", "..", "..."}
		const ClearLine = "\033[2K" // Line clear ASCII code

		for i := 0; i < totalRequestLimit; i++ {
			// Retrieved saved posts; comments
			mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, &opts)
			if err != nil {
				return nil, nil, err
			}
			fmt.Print(ClearLine)
			fmt.Printf("\r")
			fmt.Printf("Pulling saves%s", loadingDots[i%3])

			for _, post := range mySavedPosts {
				allSavedPosts = append(allSavedPosts, *post)
			}
			for _, comment := range mySavedComments {
				allSavedComments = append(allSavedComments, *comment)
			}

			// Update ListOptions.After
			opts.ListOptions.After = response.After
			time.Sleep(1 * time.Second) // Its recommend to only hit Reddit with 1 request/sec

		}

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
