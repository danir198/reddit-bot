package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http" // Add this line to import the "net/url" package
	"net/url"
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
	_, err := b.client.Comment.Upvote(b.ctx, commentID)
	return err
}

// Upvote a post
func (b *Bot) downvotePost(postID string) error {
	_, err := b.client.Post.Downvote(b.ctx, postID)
	return err
}

// Upvote a comment
func (b *Bot) downvoteComment(commentID string) error {
	_, err := b.client.Comment.Downvote(b.ctx, commentID)
	return err
}

// Follow a user
func (b *Bot) followUser(username string) error {
	path := fmt.Sprintf("api/v1/me/friends/%s", username)
	body := map[string]string{
		"action": "follow",
		"name":   username,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := b.client.NewRequest(http.MethodPut, path, url.Values{"json": {string(jsonBody)}})
	if err != nil {
		return err
	}

	_, err = b.client.Do(b.ctx, req, nil)
	return err
}

// Return selected comments based on user & body keyword
func (b *Bot) filteredComments(comments []*reddit.Comment) []*reddit.Comment {

	var c []*reddit.Comment

	for _, comment := range comments {
		if strings.Contains(strings.ToLower(comment.Body), strings.ToLower(b.upvoteBodyKeyword)) && contains(b.upvoteUserKeyword, comment.Author) {
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

	log.Printf("BOT#%v - action=%v, itemType=%v, upvoteUserKeyword=%v, upvoteBodyKeyword=%v\n", b.ID, b.action, b.itemType, b.upvoteUserKeyword, b.upvoteBodyKeyword)

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

			log.Printf("BOT#%v - Processing post ID=%v, Author=%v, Title=%v, Body=%v\n", b.ID, post.FullID, post.Author, post.Title, post.Body)

			if b.itemType == "post" {

				if contains(b.upvoteUserKeyword, post.Author) && strings.Contains(strings.ToLower(post.Body), strings.ToLower(b.upvoteBodyKeyword)) {

					hasVoted, _, err := store.HasVoted(post.FullID, b.itemType, b.ID)
					if err != nil {
						log.Printf("BOT#%v - Error checking vote status: %v\n", b.ID, err)
						continue
					}

					if hasVoted {
						log.Printf("BOT#%v - Already %s post ID=%v, Author=%v, Body=%v\n", b.ID, b.action, post.FullID, post.Author, post.Body)
					} else { // not yet voted the post
						if b.action == "upvote" {
							err = b.upvotePost(post.FullID)
							if err != nil {
								log.Printf("BOT#%v - Error upvoting post %s: %v\n", b.ID, post.FullID, err)
								continue
							}
						} else if b.action == "downvote" {
							err = b.downvotePost(post.FullID)
							if err != nil {
								log.Printf("BOT#%v - Error downvoting post %s: %v\n", b.ID, post.FullID, err)
								continue
							}
						} else { // no action
							continue
						}

						err = store.RecordVote(post.FullID, b.itemType, b.action, b.ID)
						if err != nil {
							log.Printf("BOT#%v - Error recording vote in datastore: %v\n", b.ID, err)
						}

						log.Printf("BOT#%v - Executed action=%v, itemType=%v, comment ID=%v, Author=%v, Title=%v, Body=%v\n", b.ID, b.action, b.itemType, post.ID, post.Author, post.Title, post.Body)
					}
				} else { // not matching user & content keywords
					log.Printf("BOT#%v - Skipping due to not matching keyword for users or content\n", b.ID)
				}

			} else if b.itemType == "comment" {

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

				hasVoted, _, err := store.HasVoted(randomComment.FullID, b.itemType, b.ID)
				if err != nil {
					log.Printf("BOT#%v - Error checking vote status: %v", b.ID, err)
					continue
				}

				if hasVoted {
					log.Printf("BOT#%v - Already %s comment ID=%v, Author=%v, Body=%v\n", b.ID, b.action, randomComment.FullID, randomComment.Author, randomComment.Body)
				} else { // not yet voted the comment
					if b.action == "upvote" {
						err = b.upvoteComment(randomComment.FullID)
						if err != nil {
							log.Printf("Error upvoting comment %s: %v", randomComment.FullID, err)
							continue
						}
						profile := post.Post.Author

						err = b.followUser(profile)
						if err != nil {
							log.Printf("BOT#%v - Error following profile %v: %v\n", b.ID, profile, err)
						} else {
							log.Printf("BOT#%v - Successfully followed profile %v\n", b.ID, profile)
						}

					} else if b.action == "downvote" {
						err = b.downvoteComment(randomComment.FullID)
						if err != nil {
							log.Printf("Error downvoting comment %s: %v", randomComment.FullID, err)
							continue
						}
					} else { // no action
						continue
					}

					err = store.RecordVote(randomComment.FullID, b.itemType, b.action, b.ID)
					if err != nil {
						log.Printf("BOT#%v - Error recording vote in datastore: %v", b.ID, err)
					}

					log.Printf("BOT#%v - Executed action=%v, itemType=%v, comment ID=%v, Author=%v, Body=%v\n", b.ID, b.action, b.itemType, randomComment.ID, randomComment.Author, randomComment.Body)
				}
			} else { //no itemType
				// do nothing
			}

			log.Printf("BOT#%v - Sleeping 5 seconds...\n", b.ID)
			time.Sleep(5 * time.Second)
		}

		log.Printf("BOT#%v - Sleeping 10 seconds before getting updated posts...\n", b.ID)
		time.Sleep(10 * time.Second)

	}

}
