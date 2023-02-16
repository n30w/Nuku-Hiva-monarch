package reddit

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/style"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

var (
	ResultsPerRedditRequest = 50
)

// GrabSaved reads all cached posts on the Reddit account.
// This can be used to mass refresh an entire SQL database. //TODO make this return error
func GrabSaved(postsTable, commentsTable models.DBTable, key *credentials.Key) {

	var mySavedPosts []*reddit.Post
	var mySavedComments []*reddit.Comment
	var response *reddit.Response

	// Last position of for loops
	lastPos1 := 0
	lastPos2 := 0

	totalRequests := 1

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
	client, err := reddit.NewClient(*key.RedditKey(), reddit.WithHTTPClient(httpClient))

	if err != nil {
		log.Println(style.Warn.Sprint("Login failed :("))
	} else {
		log.Println(style.Information.Sprint("Contacting Reddit API..."))
	}

	_ = style.Spinner.Start()

	for i := 0; i < totalRequests; i++ {
		mySavedPosts, mySavedComments, response, err = client.User.Saved(ctx, opts)
		if err != nil {
			log.Fatal(style.Warn.Sprint(err))
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

	log.Print(style.Result.Sprint("Saved posts and comments retrieved"))
	// log.Print(Result.Sprintf("Comments: %x", commentsTable.Rows))
	if err != nil {
		log.Fatal(style.Warn.Sprint(err))
	}
}
