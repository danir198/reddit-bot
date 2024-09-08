package main

import (
	"context"

	"github.com/vartanbeno/go-reddit/reddit"
)

type RedditClient interface {
	SubredditNewPosts(ctx context.Context, subreddit string, options *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error)
	CommentSubmit(ctx context.Context, fullname, text string) (*reddit.Comment, *reddit.Response, error)
}
