package pagedata

import (
	"fmt"
	"strings"

	"github.com/jaketreacher/gokbilgin-wiki-generator/authordata"
	"github.com/jaketreacher/gokbilgin-wiki-generator/letterdata"
)

type Page struct {
	Title string
	Text  string
}

func CreatePages(authors []*authordata.Author) []*Page {
	var letterPages []*Page
	var authorPages []*Page

	for _, author := range authors {
		letterPageMap := make(map[*letterdata.Letter]*Page)
		for _, letter := range author.Letters {
			page := createLetterPage(letter, author)

			letterPageMap[letter] = page
			letterPages = append(letterPages, page)
		}

		page := createAuthorPage(author, letterPageMap)
		authorPages = append(authorPages, page)
	}

	page := createCorrespondencePage(authorPages)

	var allPages []*Page
	allPages = append(allPages, page)
	allPages = append(allPages, authorPages...)
	allPages = append(allPages, letterPages...)

	return allPages
}

func createCorrespondencePage(authorPages []*Page) *Page {
	var authorLinks []string
	for _, page := range authorPages {
		authorLinks = append(authorLinks, createLink(urlifyTitle(page.Title), page.Title))
	}

	text := strings.Join(authorLinks, "\n\n")

	return &Page{
		Title: "Haberleşmeleri",
		Text:  text,
	}
}

func createAuthorPage(author *authordata.Author, letterPageMap map[*letterdata.Letter]*Page) *Page {
	title := fmt.Sprintf("M. Tayyib Gökbilgin'e %s Mektuplar", author.Ablative)

	var letterSections []string
	for letter, page := range letterPageMap {
		link := createLink(urlifyTitle(page.Title), page.Title)
		section := fmt.Sprintf("== %s ==\n%s", letter.Date, link)
		letterSections = append(letterSections, section)
	}

	text := strings.Join(letterSections, "\n")

	return &Page{
		title,
		text,
	}
}

func createLetterPage(letter *letterdata.Letter, author *authordata.Author) *Page {
	title := fmt.Sprintf("%s %s Tarihli Mektup", author.Ablative, letter.Date)

	var downloadLinks []string
	for _, download := range letter.Downloads {
		downloadLinks = append(downloadLinks, createLink(download.Url, download.Name))
	}

	var downloadsText string
	if len(downloadLinks) == 0 {
		downloadsText = "No Downloads"
	} else {
		downloadsText = strings.Join(downloadLinks, "\n")
	}

	text := fmt.Sprintf("%s\n== İndirilenler ==\n%s", letter.Text, downloadsText)

	return &Page{
		title,
		text,
	}
}

func urlifyTitle(title string) string {
	// TODO: Remove hard coded string
	return fmt.Sprintf("http://localhost:8080/index.php/%s", strings.ReplaceAll(title, " ", "_"))
}

func createLink(url string, text string) string {
	return fmt.Sprintf("[%s %s]", url, text)
}
