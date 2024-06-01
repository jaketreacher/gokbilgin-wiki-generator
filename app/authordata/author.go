package authordata

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jaketreacher/gokbilgin-wiki-generator/letterdata"
	"gopkg.in/yaml.v3"
)

type Author struct {
	Name      string               `yaml:"name"`
	Tags      []string             `yaml:"tags"`
	Ablative  string               `yaml:"-"`
	Letters   []*letterdata.Letter `yaml:"-"`
	Directory string               `yaml:"-"`
}

func New(directory string, letters []*letterdata.Letter) *Author {
	author := parseYaml(directory)

	author.Ablative = createAblative(author.Name)
	author.Letters = letters
	author.Directory = directory

	return author
}

func createAblative(name string) string {
	return name + "'dan"
}

func parseYaml(root string) *Author {
	path := filepath.Join(root, "author.yaml")
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatalf("+%v", err)
	}

	var author Author
	err = yaml.Unmarshal(content, &author)

	if err != nil {
		log.Fatalf("%+v", err)
	}

	return &author
}
