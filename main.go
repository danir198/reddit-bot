package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
	"log"
	"os"
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

	//
	// TODO: handle config from file or config from arguments
	//
	/*
		action := flag.String("action", "", "Action to perform: upvote or downvote (required)")
		subreddit := flag.String("subreddit", "", "Subreddit name (required)")
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
	*/

	log.SetOutput(os.Stdout)

	rbconfig := NewRedditBotConfig("config.json")
	if rbconfig != nil {
		rbconfig.Print()

		var wg sync.WaitGroup

		for _, bot := range rbconfig.Bots {

			wg.Add(1)

			b := NewBot(bot.ID, bot.Credential, []string{"aryamahzar_new", "dani198"}, "Golang", bot.Action, bot.Subreddit, bot.Actiontype)
			go b.Run(&wg)
		}

		wg.Wait()

		//fmt.Printf("Finished %sing %s in subreddit %s\n", *action, *itemType, *subreddit)
	}
}
