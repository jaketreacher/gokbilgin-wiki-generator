package main

import (
	"log"
	"os"

	"github.com/jaketreacher/gokbilgin-wiki-generator/wikiclient"
	"github.com/joho/godotenv"
)

func main() {
	username, password, endpoint := loadEnv()

	client := wikiclient.New(endpoint)

	client.Login(username, password)
	defer client.Logout()

	client.Edit("Test Page", "2 - My other content")
}

func loadEnv() (string, string, string) {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	username := os.Getenv("USER")
	password := os.Getenv("PASS")
	endpoint := os.Getenv("ENDPOINT")

	return username, password, endpoint
}
