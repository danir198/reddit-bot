package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

// Golang credentials type with JSON tagging
type Credential struct {
	REDDIT_CLIENT_ID     string `json:"REDDIT_CLIENT_ID"`
	REDDIT_CLIENT_SECRET string `json:"REDDIT_CLIENT_SECRET"`
	REDDIT_USERNAME      string `json:"REDDIT_USERNAME"`
	REDDIT_PASSWORD      string `json:"REDDIT_PASSWORD"`
}

// Config file name (JSON file)
const configFilename string = "config.json"

// List of credentials
var credentialList []Credential

// Read credentials from JSON file
func ReadCredentials() {
	//
	// Open json file
	//
	file, err := os.Open(configFilename)
	if err != nil {
		log.Printf("Error opening JSON file: %v", err)
		return
	}
	defer file.Close()

	//
	// Read json file
	//
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&credentialList); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return
	}
}

// Print out credentials variables
func PrintCredentials() {

	for idx, credential := range credentialList {
		fmt.Printf("Item number #%v\n", idx)
		credential.Print()
	}
}

// Select random item from credentials
func GetRandCredential() Credential {

	// Generate a random index in the range of the slice
	randomIndex := rand.Intn(len(credentialList))

	// Select the random item from the slice
	return credentialList[randomIndex]
}

// Print credential
func (credential Credential) Print() {
	fmt.Printf("REDDIT_CLIENT_ID: %v\n", credential.REDDIT_CLIENT_ID)
	fmt.Printf("REDDIT_CLIENT_SECRET: %v\n", credential.REDDIT_CLIENT_SECRET)
	fmt.Printf("REDDIT_USERNAME: %v\n", credential.REDDIT_USERNAME)
	fmt.Printf("REDDIT_PASSWORD: %v\n\n", credential.REDDIT_PASSWORD)
}
