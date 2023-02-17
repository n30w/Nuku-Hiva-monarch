package reddit

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/style"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var (
	ResultsPerRedditRequest = 50 // Number of saved objects to return
	totalRequests           = 1  // Amount of requests to make to Reddit.
)

// GrabSaved reads all cached posts on the Reddit account.
// This can be used to mass refresh an entire SQL database.
func GrabSaved(postsTable, commentsTable models.DBTable, key credentials.Authenticator) error {

	var mySavedPosts []*reddit.Post
	var mySavedComments []*reddit.Comment
	var response *reddit.Response

	// Last position of for loops
	lastPos1, lastPos2 := 0, 0

	ctx := context.Background()
	opts := &reddit.ListUserOverviewOptions{
		ListOptions: reddit.ListOptions{
			Limit:  ResultsPerRedditRequest,
			After:  "",
			Before: "",
		},
		Sort: "new",
		Time: "all",
	}

	// Establish connection to Reddit API
	httpClient := &http.Client{Timeout: time.Second * 30}
	redditKey := key.Use().(*reddit.Credentials)
	client, err := reddit.NewClient(*redditKey, reddit.WithHTTPClient(httpClient))

	if err != nil {
		return errors.New("authentication with Reddit API failed:" + err.Error())
	} else {
		log.Println(style.Information.Sprint("Contacting Reddit API..."))
	}

	_ = style.Spinner.Start()

	for i := 0; i < totalRequests; i++ {
		mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, opts)
		if err != nil {
			return err
		}

		// TODO go routine optimization can occur here
		for _, post := range mySavedPosts {
			postsTable.Rows[lastPos1] = models.NewRow(
				0,
				post.Title,
				post.Permalink,
				post.SubredditName,
				post.URL,
			)
			lastPos1++
		}

		// TODO go routine optimization can occur here
		for _, comment := range mySavedComments {
			commentsTable.Rows[lastPos2] = models.NewRow(
				0,
				comment.Author,
				comment.Body,
				comment.Permalink,
				comment.SubredditName,
			)
			lastPos2++
		}

		style.Spinner.Message(opts.ListOptions.After)

		opts.ListOptions.After = response.After
		time.Sleep(1 * time.Second) // Its recommend to hit Reddit with only 1 request/sec
	}

	_ = style.Spinner.Stop()

	log.Print(style.Result.Sprint("saved posts and comments retrieved"))
	// log.Print(Result.Sprintf("Comments: %x", commentsTable.Rows))
	if err != nil {
		return err
	}

	return nil
}
