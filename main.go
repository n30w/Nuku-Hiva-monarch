package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	/*
		Establish connection
	*/
	httpClient := &http.Client{Timeout: time.Second * 30}
	client, _ := reddit.NewClient(credentials, reddit.WithHTTPClient(httpClient))

	/*
		Options to satisfy client.User.Saved's 2nd parameter
	*/
	opts := reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit:  100,
			After:  "",
			Before: "",
		},
		Sort: "new",
		Time: "all",
	}

	/*
		Retrieve saved posts, store in var
	*/
	mySavedPosts, _, _, err := client.User.Saved(ctx, &opts)
	if err != nil {
		return
	}

	/*
		Retrieve and print saved posts by time
	*/
	for _, post := range mySavedPosts {
		fmt.Printf("%s | %s\n", post.Title, post.URL)
	}

}
