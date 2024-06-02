package page

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/author"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/letter"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type Page struct {
	Title string
	Text  string
}

type LetterSection struct {
	Date    string
	Content string
}

func CreatePages(authors []*author.Author) []*Page {
	var letterPages []*Page
	var authorPages []*Page

	authorPageMap := make(map[*author.Author]*Page)
	for _, author := range authors {
		letterPageMap := make(map[*letter.Letter]*Page)
		for _, letter := range author.Letters {
			page := createLetterPage(letter, author)

			letterPageMap[letter] = page
			letterPages = append(letterPages, page)
		}

		page := createAuthorPage(author, letterPageMap)
		authorPages = append(authorPages, page)
		authorPageMap[author] = page
	}

	page := createCorrespondencePage(authorPageMap)

	var allPages []*Page
	allPages = append(allPages, page)
	allPages = append(allPages, authorPages...)
	allPages = append(allPages, letterPages...)

	return allPages
}

func createCorrespondencePage(authorPageMap map[*author.Author]*Page) *Page {
	authors := make([]*author.Author, 0, len(authorPageMap))
	for key := range authorPageMap {
		authors = append(authors, key)
	}
	sort.Slice(authors, func(i, j int) bool {
		collator := collate.New(language.Turkish)
		return collator.CompareString(authors[i].SortKey, authors[j].SortKey) < 0
	})

	var authorLinks []string
	for _, author := range authors {
		page := authorPageMap[author]
		authorLinks = append(authorLinks, createInternalLink(page.Title))
	}

	text := strings.Join(authorLinks, "\n\n")

	return &Page{
		Title: "Haberleşmeleri",
		Text:  text,
	}
}

func createAuthorPage(author *author.Author, letterPageMap map[*letter.Letter]*Page) *Page {
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

func createLetterPage(letter *letter.Letter, author *author.Author) *Page {
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
