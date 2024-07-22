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

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Replace with your Reddit app credentials
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
	repliedPost := make(map[string]struct{})

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
			// Check if the post has already been replied to
			if _, replied := repliedPost[post.ID]; replied {
				continue

			}
			replyMessage := "Hello! This is an automated reply."
			_, _, err := client.Comment.Submit(ctx, post.FullID, replyMessage)
			if err != nil {
				log.Printf("Error replying to post %s: %v", post.ID, err)
				continue
			}

			fmt.Printf("Replied to post: %s\n", post.Title)
			time.Sleep(5 * time.Second)
		}

		// Sleep for a while before checking again
		time.Sleep(10 * time.Minute)
	}
}
