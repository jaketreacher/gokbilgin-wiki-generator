package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jaketreacher/gokbilgin-wiki-generator/authordata"
	"github.com/jaketreacher/gokbilgin-wiki-generator/letterdata"
	"github.com/jaketreacher/gokbilgin-wiki-generator/pagedata"
	"github.com/jaketreacher/gokbilgin-wiki-generator/wikiclient"
	"github.com/joho/godotenv"
)

func main() {
	username, password, endpoint := loadEnv()
	client := wikiclient.New(endpoint)
	client.Login(username, password)
	defer client.Logout()

	client.UserInfoQuery()
	input := "../data/letters"

	author_dirs := getDirs(input)
	var authors []*authordata.Author
	for _, author_dir := range author_dirs {
		letter_dirs := getDirs(author_dir)

		var letters []*letterdata.Letter
		for _, letter_dir := range letter_dirs {
			letter := letterdata.New(letter_dir)
			letters = append(letters, letter)
		}

		author := authordata.New(author_dir, letters)
		authors = append(authors, author)
	}

	pages := pagedata.CreatePages(authors)

	for _, page := range pages {
		client.Edit(page.Title, page.Text)
	}
}

func loadEnv() (string, string, string) {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("%+v", err)
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	endpoint := os.Getenv("ENDPOINT")

	return username, password, endpoint
}

func getDirs(root string) []string {
	entries, err := os.ReadDir(root)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(root, entry.Name())
			dirs = append(dirs, path)
		}
	}

	return dirs
}
