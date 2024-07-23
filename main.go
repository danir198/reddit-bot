package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Load environment variables from .env file if present
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}
}

// Check if a post has already been replied to
func hasReplied(postID string, repliedPosts map[string]bool) bool {
	_, exists := repliedPosts[postID]
	return exists
}

// Mark a post as replied to
func markReplied(postID string, repliedPosts map[string]bool) {
	repliedPosts[postID] = true
}

func main() {
	loadEnv()

	// Replace with your Reddit app credentials from environment variables
	credentials := reddit.Credentials{
		ID:       os.Getenv("REDDIT_CLIENT_ID"),
		Secret:   os.Getenv("REDDIT_CLIENT_SECRET"),
		Username: os.Getenv("REDDIT_USERNAME"),
		Password: os.Getenv("REDDIT_PASSWORD"),
	}

	// Create a new Reddit client
	client, err := reddit.NewClient(credentials)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	subreddit := "test_learning_bot_gol" // Replace with your subreddit of choice

	// Keep track of replied posts to avoid duplicates
	repliedPosts := make(map[string]bool)

	for {
		// Fetch new posts
		posts, _, err := client.Subreddit.NewPosts(ctx, subreddit, &reddit.ListOptions{
			Limit: 10,
		})
		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			continue
		}

		for _, post := range posts {
			if hasReplied(post.FullID, repliedPosts) {
				continue
			}

			replyMessage := "Hello! This is an automated reply."
			_, _, err := client.Comment.Submit(ctx, post.FullID, replyMessage)
			if err != nil {
				log.Printf("Error replying to post %s: %v", post.ID, err)
				continue
			}

			markReplied(post.FullID, repliedPosts)
			fmt.Printf("Replied to post: %s\n", post.Title)

			// Sleep for a while before replying to the next post to respect rate limits
			time.Sleep(5 * time.Second)
		}

		// Sleep for a while before checking for new posts again
		time.Sleep(1 * time.Minute)
	}
}
