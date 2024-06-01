package pagedata

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/jaketreacher/gokbilgin-wiki-generator/authordata"
	"github.com/jaketreacher/gokbilgin-wiki-generator/letterdata"
)

type Page struct {
	Title string
	Text  string
}

type LetterSection struct {
	Date    string
	Content string
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
		authorLinks = append(authorLinks, createInternalLink(page.Title))
	}

	slices.Sort(authorLinks)

	text := strings.Join(authorLinks, "\n\n")

	return &Page{
		Title: "Haberleşmeleri",
		Text:  text,
	}
}

func createAuthorPage(author *authordata.Author, letterPageMap map[*letterdata.Letter]*Page) *Page {
	title := fmt.Sprintf("M. Tayyib Gökbilgin'e %s Mektuplar", author.Ablative)

	var letterSections []*LetterSection
	for letter, page := range letterPageMap {
		link := createInternalLink(page.Title)
		content := fmt.Sprintf("== %s ==\n%s", letter.Date, link)
		letterSections = append(letterSections, &LetterSection{Date: letter.Date, Content: content})
	}

	sort.Slice(letterSections, func(i, j int) bool {
		return letterSections[i].Date < letterSections[j].Date
	})

	var letterSectionTexts []string
	for _, section := range letterSections {
		letterSectionTexts = append(letterSectionTexts, section.Content)
	}
	text := strings.Join(letterSectionTexts, "\n\n")

	return &Page{
		title,
		text,
	}
}

func createLetterPage(letter *letterdata.Letter, author *authordata.Author) *Page {
	title := fmt.Sprintf("%s %s Tarihli Mektup", author.Ablative, letter.Date)

	var downloadLinks []string
	for _, download := range letter.Downloads {
		downloadLinks = append(downloadLinks, createExternalLink(download.Url, download.Name))
	}

	var downloadsText string
	if len(downloadLinks) == 0 {
		downloadsText = "No Downloads"
	} else {
		downloadsText = strings.Join(downloadLinks, "\n\n")
	}

	text := fmt.Sprintf("%s\n== İndirilenler ==\n%s", letter.Text, downloadsText)

	return &Page{
		title,
		text,
	}
}

func createInternalLink(title string) string {
	return fmt.Sprintf("[[%s]]", title)
}

func createExternalLink(url string, text string) string {
	return fmt.Sprintf("[%s %s]", url, text)
}
