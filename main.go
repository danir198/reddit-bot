package main

//
// Reddit reference:
// - https://www.reddit.com/dev/api/
// - https://www.reddit.com/prefs/apps
//

import (
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

	log.SetOutput(os.Stdout)

	var wg sync.WaitGroup
	ReadCredentials()

	for idx, credential := range CredentialList {

		wg.Add(1)

		b := NewBot(strconv.Itoa(idx), credential, []string{"aryamahzar_new", "dani198"}, "Golang")
		go b.Run(&wg)
	}

	wg.Wait()
}
