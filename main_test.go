package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type MockRedditClient struct{}

func (m *MockRedditClient) SubredditNewPosts(ctx context.Context, subreddit string, options *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error) {
	posts := []*reddit.Post{
		{FullID: "t3_post1", Title: "Post 1"},
		{FullID: "t3_post2", Title: "Post 2"},
	}
	return posts, nil, nil
}

func (m *MockRedditClient) CommentSubmit(ctx context.Context, fullname, text string) (*reddit.Comment, *reddit.Response, error) {
	return &reddit.Comment{FullID: fullname}, nil, nil
}

func TestHasReplied(t *testing.T) {
	repliedPosts := map[string]bool{
		"t3_post1": true,
	}

	assert.True(t, hasReplied("t3_post1", repliedPosts))
	assert.False(t, hasReplied("t3_post2", repliedPosts))
}

func TestMarkReplied(t *testing.T) {
	repliedPosts := make(map[string]bool)
	markReplied("t3_post1", repliedPosts)
	assert.True(t, repliedPosts["t3_post1"])
}

// func TestLoadEnv(t *testing.T) {
// 	loadEnv()
// 	assert.Equal(t, "your_client_id", os.Getenv("REDDIT_CLIENT_ID"))
// 	assert.Equal(t, "your_client_secret", os.Getenv("REDDIT_CLIENT_SECRET"))
// 	assert.Equal(t, "your_reddit_username", os.Getenv("REDDIT_USERNAME"))
// 	assert.Equal(t, "your_reddit_password", os.Getenv("REDDIT_PASSWORD"))
// }

func TestRedditBot(t *testing.T) {
	client := &MockRedditClient{}
	repliedPosts := make(map[string]bool)
	ctx := context.Background()
	subreddit := "test_learning_bot_gol"

	posts, _, err := client.SubredditNewPosts(ctx, subreddit, &reddit.ListOptions{Limit: 10})
	assert.NoError(t, err)

	for _, post := range posts {
		if hasReplied(post.FullID, repliedPosts) {
			continue
		}

		replyMessage := "Hello! This is an automated reply."
		_, _, err := client.CommentSubmit(ctx, post.FullID, replyMessage)
		assert.NoError(t, err)

		markReplied(post.FullID, repliedPosts)
		assert.True(t, hasReplied(post.FullID, repliedPosts))
		time.Sleep(5 * time.Second)
	}
}
