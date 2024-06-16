package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/author"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/letter"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/page"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/wikiclient"
	"github.com/joho/godotenv"
)

func main() {
	input := os.Args[1]

	info, err := os.Stat(input)
	if os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("%v is invalid", input)
	}

	username, password, endpoint := loadEnv()
	client := wikiclient.New(endpoint)
	client.Login(username, password)
	defer client.Logout()

	client.UserInfoQuery()

	authorDirs := getDirs(input)
	var authors []*author.Author
	for _, authorDir := range authorDirs {
		lettersDir := getDirs(authorDir)

		var letters []*letter.Letter
		for _, letter_dir := range lettersDir {
			letter := letter.New(letter_dir)
			letters = append(letters, letter)
		}

		author := author.New(authorDir, letters)
		authors = append(authors, author)
	}

	pages := page.CreatePages(authors)

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
