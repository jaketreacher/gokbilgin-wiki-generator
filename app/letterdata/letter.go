package letterdata

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Letter struct {
	Date      string `yaml:"date"`
	Downloads []struct {
		Name string `yaml:"name"`
		Url  string `yaml:"url"`
	} `yaml:"downloads"`
	Text      string `yaml:"-"`
	Directory string `yaml:"-"`
}

func New(directory string) Letter {
	letter := parseYaml(directory)
	text := readTranslations(directory)

	letter.Text = text
	letter.Directory = directory

	return letter
}

func parseYaml(root string) Letter {
	path := filepath.Join(root, "letter.yaml")
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("+%v", err)
	}

	var letter Letter
	err = yaml.Unmarshal(content, &letter)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	return letter
}

func readTranslations(root string) string {
	path := filepath.Join(root, "translation.tur.txt")
	content, err := os.ReadFile(path)

	if err != nil {
		return "No translations found"
	} else {
		return string(content)
	}
}
