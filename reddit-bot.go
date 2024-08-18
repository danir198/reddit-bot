package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"reddit-bot/datastore"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Bot represents a Reddit bot instance.
type Bot struct {
	ID           string
	ctx          context.Context
	credential   botCredentialType
	client       *reddit.Client
	subreddit    string
	repliedPosts map[string]bool

	upvoteUserKeyword []string
	upvoteBodyKeyword string
	action            string
	itemType          string
}

// NewBot creates a new Bot with a specific account.
func NewBot(ID string, c botCredentialType, upvoteUserKeyword []string, upvoteBodyKeyword string, action, subreddit, itemType string) *Bot {

	b := Bot{
		ID:                ID,
		credential:        c,
		upvoteUserKeyword: upvoteUserKeyword,
		upvoteBodyKeyword: upvoteBodyKeyword,
		action:            action,
		subreddit:         subreddit,
		itemType:          itemType,
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
	_, err := b.client.Comment.Downvote(b.ctx, commentID)
	return err
}

// Upvote a post
func (b *Bot) downvotePost(postID string) error {
	_, err := b.client.Post.Upvote(b.ctx, postID)
	return err
}

// Upvote a comment
func (b *Bot) downvoteComment(commentID string) error {
	_, err := b.client.Comment.Downvote(b.ctx, commentID)
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

var store datastore.VoteStore

func (b Bot) Run(wg *sync.WaitGroup) {

	defer wg.Done()

	b.ctx = context.Background()

	log.Printf("Running BOT#%v...", b.ID)
	log.Printf("BOT#%v - credential: %v", b.ID, b.credential)

	var err error
	store, err = datastore.NewSQLiteStore("reddit_bot.db")
	if err != nil {
		log.Fatalf("Failed to create datastore: %v", err)
	}
	defer store.Close()
	defer wg.Done()

	for {

		// Fetch new posts
		log.Printf("BOT#%v - Fetching new posts\n", b.ID)
		posts, _, err := b.client.Subreddit.NewPosts(b.ctx, b.subreddit, &reddit.ListOptions{
			Limit: 50,
		})

		if err != nil {
			log.Printf("BOT#%v - Error fetching posts: %v\n", b.ID, err)
			log.Printf("Error fetching posts: %v", err)
			continue
		}

		// Process each new posts
		for _, post := range posts {

			log.Printf("BOT#%v - Processing post ID=%v, Author=%v, Title=%v\n", b.ID, post.FullID, post.Author, post.Title)

			post, _, err := b.client.Post.Get(b.ctx, post.ID)
			if err != nil {
				log.Printf("BOT#%v - Error fetching comments for post ID=%v: %v", b.ID, post.Post.ID, err)
				return
			}

			filteredComments := b.filteredComments(post.Comments)
			if len(filteredComments) == 0 {
				log.Printf("BOT#%v - No filtered comments for post ID=%v", b.ID, post.Post.ID)
				continue
			}

			randomComment := filteredComments[rand.Intn(len(filteredComments))]
			log.Printf("BOT#%v - Selected random comment ID=%v", b.ID, randomComment.FullID)

			hasVoted, prevAction, err := store.HasVoted(randomComment.FullID, b.ID)
			if err != nil {
				log.Printf("BOT#%v - Error checking vote status: %v", b.ID, err)
				continue
			}

			if hasVoted {
				log.Printf("BOT#%v - Already %s comment ID=%v, Author=%v, Body=%v\n", b.ID, prevAction, randomComment.ID, randomComment.Author, randomComment.Body)
				continue
			}

			if b.action == "upvote" && b.itemType == "comment" {
				err = b.upvoteComment(randomComment.FullID)
				if err != nil {
					log.Printf("Error upvoting comment %s: %v", randomComment.FullID, err)
					continue
				}
			} else if b.action == "downvote" && b.itemType == "comment" {

				err = b.downvotePost(randomComment.FullID)
				if err != nil {
					log.Printf("Error downvoting comment %s: %v", randomComment.FullID, err)
					continue
				}
			}

			err = store.RecordVote(randomComment.FullID, "comment", b.action, b.ID)
			if err != nil {
				log.Printf("BOT#%v - Error recording vote in datastore: %v", b.ID, err)
			}

			fmt.Printf("action=%v, comment ID=%v, Author=%v, Body=%v: %v\n", b.action, b.ID, randomComment.ID, randomComment.Author, randomComment.Body)

			log.Printf("BOT#%v - Sleeping 5 seconds...\n", b.ID)
			time.Sleep(5 * time.Second)
		}

		log.Printf("BOT#%v - Sleeping 10 seconds...\n", b.ID)
		time.Sleep(10 * time.Second)

	}

}
