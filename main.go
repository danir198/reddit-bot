package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

// Load environment variables from .env file if present
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}
}

func main() {
	action := flag.String("action", "", "Action to perform: upvote or downvote (required)")
	subreddit := flag.String("subreddit", "", "Subreddit name (required)")
	itemID := flag.String("id", "", "ID of the post or comment to vote on (required)")
	itemType := flag.String("type", "", "Type of item to vote on: post or comment (required)")

	flag.Parse()

	if *action == "" || *subreddit == "" || *itemType == "" {
		fmt.Println("All flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *action != "upvote" && *action != "downvote" {
		fmt.Println("Invalid action. Use 'upvote' or 'downvote'.")
		os.Exit(1)
	}

	if *itemType != "post" && *itemType != "comment" {
		fmt.Println("Invalid item type. Use 'post' or 'comment'.")
		os.Exit(1)
	}

	log.SetOutput(os.Stdout)
	ReadCredentials()

	var wg sync.WaitGroup

	for idx, credential := range CredentialList {

		wg.Add(1)

		b := NewBot(strconv.Itoa(idx), credential, []string{"aryamahzar_new", "dani198"}, "Golang", *action, *subreddit, *itemType)
		go b.Run(&wg)
	}

	wg.Wait()
	fmt.Printf("Finished %sing %s %s in subreddit %s\n", *action, *itemType, *itemID, *subreddit)
}
