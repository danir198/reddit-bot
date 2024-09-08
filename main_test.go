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
	return []*reddit.Post{
		{FullID: "t3_post1"},
		{FullID: "t3_post2"},
	}, nil, nil
}

func (m *MockRedditClient) CommentSubmit(ctx context.Context, fullname, text string) (*reddit.Comment, *reddit.Response, error) {
	return &reddit.Comment{ID: "t1_comment1"}, nil, nil
}

func TestHasReplied(t *testing.T) {
	bot := &Bot{
		repliedPosts: make(map[string]bool),
	}

	assert.False(t, bot.hasReplied("t3_post1"))
	bot.markReplied("t3_post1")
	assert.True(t, bot.hasReplied("t3_post1"))
}

func TestMarkReplied(t *testing.T) {
	bot := &Bot{
		repliedPosts: make(map[string]bool),
	}

	bot.markReplied("t3_post1")
	assert.True(t, bot.repliedPosts["t3_post1"])
}

func TestRedditBot(t *testing.T) {
	client := &MockRedditClient{}
	bot := &Bot{
		client:       client,
		repliedPosts: make(map[string]bool),
	}
	ctx := context.Background()
	subreddit := "test_learning_bot_gol"

	posts, _, err := client.SubredditNewPosts(ctx, subreddit, &reddit.ListOptions{Limit: 10})
	assert.NoError(t, err)

	for _, post := range posts {
		if bot.hasReplied(post.FullID) {
			continue
		}

		replyMessage := "Hello! This is an automated reply."
		_, _, err := client.CommentSubmit(ctx, post.FullID, replyMessage)
		assert.NoError(t, err)

		bot.markReplied(post.FullID)
		assert.True(t, bot.hasReplied(post.FullID))
		time.Sleep(5 * time.Second)
	}
}
