package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/avast/retry-go"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/author"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/letter"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/page"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/wikiclient"
	"github.com/joho/godotenv"
)

const BatchSize = 10

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

	log.Printf("Processing %d page(s)\n", len(pages))

	var wg sync.WaitGroup
	inputQueue := make(chan *page.Page)

	go func() {
		defer close(inputQueue)
		for _, page := range pages {
			inputQueue <- page
		}
	}()

	for i := 0; i < BatchSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range inputQueue {
				err := retry.Do(func() error {
					result, err := client.Edit(page.Title, page.Text)

					if err != nil {
						return err
					}

					log.Printf("Page processed: %s", result.Edit.Title)
					return nil
				})
				if err != nil {
					log.Println(err)
				}
			}
		}()
	}

	wg.Wait()
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
