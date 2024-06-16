package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DateRegexFull = `\d{4}\.\d{2}\.\d{2}`
	DateRegexYear = `\d{4}`
	UnknownAuthor = "Unorganized"
	UnknownDate   = "0000.00.00"
)

type DocumentType string

var Document = struct {
	Original    DocumentType
	Translation DocumentType
}{
	Original:    "original",
	Translation: "tercume",
}

func main() {
	input := os.Args[1]

	info, err := os.Stat(input)
	if os.IsNotExist(err) || !info.IsDir() {
		log.Fatalf("%v is invalid", input)
	}

	root := input

	baseDepth := strings.Count(root, string(os.PathSeparator))
	filepath.WalkDir(root, func(path string, d fs.DirEntry, _ error) error {
		if isUnsupportedFile(d) {
			return nil
		}

		currentDepth := strings.Count(path, string(os.PathSeparator))
		relativeDepth := currentDepth - baseDepth

		filename := filepath.Base(path)
		filedate := extractDate(path, relativeDepth)
		if filedate == "" {
			filedate = UnknownDate
		}

		authorDir, err := getAuthorDir(path, relativeDepth)
		if err != nil {
			msg := map[string]string{
				"path": path,
				"err":  err.Error(),
			}
			log.Printf("Skipping file due to error: %+v\n", msg)
			return nil
		}

		documentType, err := getDocumentType(filename)
		if err != nil {
			msg := map[string]string{
				"filename": filename,
				"err":      err.Error(),
			}
			log.Printf("Skipping file due to error: %+v\n", msg)
		}

		newName := fmt.Sprintf("%s %s %s%s", authorDir, filedate, documentType, filepath.Ext(filename))
		newPath := filepath.Join(root, authorDir, filedate, newName)

		if path != newPath {
			os.MkdirAll(filepath.Dir(newPath), 0777)
			err := os.Rename(path, newPath)
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("MOVED\t%s => %s", path, newPath)
		}

		return nil
	})
}

func isUnsupportedFile(d fs.DirEntry) bool {
	return d.IsDir() || d.Name()[0] == '.' || strings.ToLower(filepath.Ext(d.Name())) != ".pdf"
}

func extractDate(path string, depth int) string {
	re := regexp.MustCompile(DateRegexFull)

	if depth == 3 {
		parentDir := filepath.Base(filepath.Dir(path))
		return re.FindString(parentDir)
	}

	filename := filepath.Base(path)

	fullDate := re.FindString(filename)
	if fullDate != "" {
		return fullDate
	}

	yearOnlyDate := regexp.MustCompile(DateRegexYear).FindString(filename)
	if yearOnlyDate != "" {
		return fmt.Sprintf("%s.00.00", yearOnlyDate)
	}

	return ""
}

// Determine the author directory from a given file path based on the specified depth.
// The directory structure is structured such that the author directories is one level
// down from the root.
//
// Depth examples:
// - 1: ./root/file.pdf
// - 2: ./root/author/file.pdf
// - 3: ./root/author/date/file.pdf
func getAuthorDir(path string, depth int) (string, error) {
	switch depth {
	case 1:
		return UnknownAuthor, nil
	case 2:
		return filepath.Base(filepath.Dir(path)), nil
	case 3:
		return filepath.Base(filepath.Dir(filepath.Dir(path))), nil
	default:
		return "", errors.New("unsupported depth")
	}
}

func getDocumentType(name string) (DocumentType, error) {
	name = strings.ToLower(name)

	if strings.Contains(name, string(Document.Original)) {
		return Document.Original, nil
	}

	if strings.Contains(name, string(Document.Translation)) {
		return Document.Translation, nil
	}

	return "", errors.New("unsupported document")
}
