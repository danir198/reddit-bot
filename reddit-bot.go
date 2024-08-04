package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
	"context"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Bot represents a Reddit bot instance.
type Bot struct {
	ID           string
	ctx          context.Context
	credential   Credential
	client       *reddit.Client
	subreddit    string
	repliedPosts map[string]bool

	upvoteUserKeyword []string
	upvoteBodyKeyword string
}

// NewBot creates a new Bot with a specific account.
func NewBot(ID string, c Credential, upvoteUserKeyword []string, upvoteBodyKeyword string) *Bot {

	b := Bot{
		ID:                ID,
		credential:        c,
		upvoteUserKeyword: upvoteUserKeyword,
		upvoteBodyKeyword: upvoteBodyKeyword,
	}

	var err error

	// Create a new Reddit client
	b.client, err = reddit.NewClient(reddit.Credentials{
		ID:       c.REDDIT_CLIENT_ID,
		Secret:   c.REDDIT_CLIENT_SECRET,
		Username: c.REDDIT_USERNAME,
		Password: c.REDDIT_PASSWORD,
	})

	if err != nil {
		log.Fatal(err)
	}

	b.repliedPosts = make(map[string]bool)

	return &b
}

// Check if a post has already been replied to
func (b *Bot) hasReplied(postID string) bool {
	return b.repliedPosts[postID]
}

// Mark a post as replied to
func (b *Bot) markReplied(postID string) {
	b.repliedPosts[postID] = true
}

// Upvote a post
func (b *Bot) upvotePost(postID string) error {
	_, err := b.client.Post.Upvote(b.ctx, postID)
	return err
}

// Upvote a comment
func (b *Bot) upvoteComment(commentID string) error {
	_, err := b.client.Comment.Upvote(b.ctx, commentID)
	return err
}

// Return selected comments based on user & body keyword
func (b *Bot) filteredComments(comments []*reddit.Comment) []*reddit.Comment {

	var c []*reddit.Comment

	for _, comment := range comments {
		if strings.Contains(comment.Body, b.upvoteBodyKeyword) && contains(b.upvoteUserKeyword, comment.Author) {
			c = append(c, comment)
		}
	}
	return c
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func (b Bot) Run(wg *sync.WaitGroup) {

	defer wg.Done()

	b.ctx = context.Background()
	b.subreddit = `test_learning_bot_gol`
	//b.subreddit = "my_subreddit_test"

	log.Printf("Running BOT#%v...", b.ID)
	log.Printf("BOT#%v - credential: %v", b.ID, b.credential)

	for {

		// Fetch new posts
		log.Printf("BOT#%v - Fetching new posts\n", b.ID)
		posts, _, err := b.client.Subreddit.NewPosts(b.ctx, b.subreddit, &reddit.ListOptions{
			Limit: 10,
		})

		if err != nil {
			log.Printf("BOT#%v - Error fetching posts: %v\n", b.ID, err)
			log.Printf("Error fetching posts: %v", err)
			continue
		}

		// Process each new posts
		for _, post := range posts {

			log.Printf("BOT#%v - Processing post ID=%v, Author=%v, Title=%v\n", b.ID, post.FullID, post.Author, post.Title)

			/*
				//
				// Process reply post
				//
				if b.hasReplied(post.FullID) {
					continue
				}

				replyMessage := "Hello! This is an automated reply."
				_, _, err := b.client.Comment.Submit(b.ctx, post.FullID, replyMessage)
				if err != nil {
					log.Printf("BOT#%v - Error replying to post ID=%: %v\n", b.ID, post.ID, err)
					continue
				}

				log.Printf("BOT#%v - Replied to post ID=%v, Author=%v, Title=%v\n", b.ID, post.ID, post.Title)
				b.markReplied(post.FullID)
			*/

			//
			// Process upvote post
			//

			post, _, err := b.client.Post.Get(b.ctx, post.ID)
			if err != nil {
				log.Printf("BOT#%v - Error fetching comments for post ID=%v: %v", b.ID, post.Post.ID, err)
				return
			}

			filteredComments := b.filteredComments(post.Comments)
			if len(filteredComments) > 0 {
				randomComment := filteredComments[rand.Intn(len(filteredComments))]

				if randomComment.Likes != nil && *(randomComment.Likes) {
					log.Printf("BOT#%v - Already Upvoted comment ID=%v, Author=%v, Body=%v\n", b.ID, randomComment.ID, randomComment.Author, randomComment.Body)
				} else {

					err = b.upvoteComment(randomComment.FullID)
					if err != nil {
						log.Printf("BOT#%v - Error upvoting comment ID=%v, Author=%v, Body=%v: %v\n", b.ID, randomComment.ID, randomComment.Author, randomComment.Body, err)
					} else {
						log.Printf("BOT#%v - Upvoted comment ID=%v, Author=%v, Body=%v\n", b.ID, randomComment.ID, randomComment.Author, randomComment.Body)
					}
				}
			}

			log.Printf("BOT#%v - Sleeping 5 seconds...\n", b.ID)
			time.Sleep(5 * time.Second)
		}

		log.Printf("BOT#%v - Sleeping 10 seconds...\n", b.ID)
		time.Sleep(10 * time.Second)

	}
}
