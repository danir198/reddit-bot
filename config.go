package main

import (
	"encoding/json"
	"log"
	"os"
)

// Reddit-bot configuration
type botCredentialType struct {
	REDDIT_CLIENT_ID     string `json:"REDDIT_CLIENT_ID"`
	REDDIT_CLIENT_SECRET string `json:"REDDIT_CLIENT_SECRET"`
	REDDIT_USERNAME      string `json:"REDDIT_USERNAME"`
	REDDIT_PASSWORD      string `json:"REDDIT_PASSWORD"`
}

// TODO: add upvoteUserKeyword and upvoteBodyKeyword in botConfigType
type botConfigType struct {
	ID                string            `json:"ID"`
	Subreddit         string            `json:"subreddit"`
	Action            string            `json:"action"`     // upvote or downvote
	Actiontype        string            `json:"actiontype"` // post or comment
	UpvoteUserKeyword []string          `json:"upvoteUserKeyword"`
	UpvoteBodyKeyword string            `json:"upvoteBodyKeyword"`
	Credential        botCredentialType `json:"credential"`
}

type redditBotConfigType struct {
	App_name string          `json:"app_name"`
	Version  string          `json:"version"`
	Bots     []botConfigType `json:"bots"`
}

// Initialize reddit-bot config, read from json config file
func NewRedditBotConfig(configFilename string) *redditBotConfigType {

	var RedditBotConfig redditBotConfigType

	//
	// Open json file
	//
	log.Printf("Reading reddit-bot configuration file: %v\n", configFilename)
	file, err := os.Open(configFilename)
	if err != nil {
		log.Printf("Error opening JSON file: %v", err)
		return nil
	}
	defer file.Close()

	//
	// Read json file
	//
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&RedditBotConfig)

	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil
	} else {
		return &RedditBotConfig
	}
}

// Print config
func (r redditBotConfigType) Print() {
	log.Printf("reddit-bot config: %v\n", r)
}

/*
func main() {
	rbconfig := NewRedditBotConfig("config.json")
	if rbconfig != nil {
		rbconfig.Print()
	}
}
*/
