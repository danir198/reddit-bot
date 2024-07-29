package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/vartanbeno/go-reddit/v2/reddit"
	"golang.org/x/exp/rand"
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

// Upvote a post
func upvotePost(client *reddit.Client, ctx context.Context, postID string) error {
	_, err := client.Post.Upvote(ctx, postID)
	return err
}

func upvoteComment(client *reddit.Client, ctx context.Context, commentID string) error {
	_, err := client.Comment.Upvote(ctx, commentID)
	return err
}

func filterComments(comments []*reddit.Comment, keyword string, usernames []string) []*reddit.Comment {
	var filteredComments []*reddit.Comment
	for _, comment := range comments {
		if strings.Contains(comment.Body, keyword) && contains(usernames, comment.Author) {
			filteredComments = append(filteredComments, comment)
			fmt.Println(comment)
		}
	}
	return filteredComments
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
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

	// Check REDDIT env vairables
	fmt.Println("REDDIT_CLIENT_ID = ", credentials.ID)
	fmt.Println("REDDIT_CLIENT_SECRET = ", credentials.Secret)
	fmt.Println("REDDIT_USERNAME = ", credentials.Username)
	fmt.Println("REDDIT_PASSWORD = ", credentials.Password)

	// Create a new Reddit client
	client, err := reddit.NewClient(credentials)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	//subreddit := `test_learning_bot_gol` // Replace with your subreddit of choice
	subreddit := "my_subreddit_test" // Replace with your subreddit of choice

	// Keep track of replied posts to avoid duplicates
	repliedPosts := make(map[string]bool)
	keyword := "Golang"
	usernames := []string{"dani198", "specificUser2"}

	for {
		// Fetch new posts
		fmt.Println("Fetching new post...")
		posts, _, err := client.Subreddit.NewPosts(ctx, subreddit, &reddit.ListOptions{
			Limit: 10,
		})

		if err != nil {
			log.Printf("Error fetching posts: %v", err)
			continue
		}

		for _, post := range posts {

			fmt.Printf("%v - %v - %v\n", post.Author, post.Title, post.Body)

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

			post, _, err := client.Post.Get(ctx, post.ID)
			if err != nil {
				log.Printf("Error fetching post %s: %v", post.Post.ID, err)
				return
			}

			filteredComments := filterComments(post.Comments, keyword, usernames)
			if len(filteredComments) > 0 {
				randomComment := filteredComments[rand.Intn(len(filteredComments))]
				err = upvoteComment(client, ctx, randomComment.FullID)
				if err != nil {
					log.Printf("Error upvoting comment %s: %v", randomComment.ID, err)
				} else {
					fmt.Printf("Upvoted comment by %s: %s\n", randomComment.Author, randomComment.Body)
				}
			}

			time.Sleep(5 * time.Second)
		}
		time.Sleep(1 * time.Minute)

	}
}
