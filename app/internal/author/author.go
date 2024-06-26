package author

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/letter"
	"github.com/jaketreacher/gokbilgin-wiki-generator/internal/turkishsuffix"
	"gopkg.in/yaml.v3"
)

type Author struct {
	Name      string           `yaml:"name"`
	Tags      []string         `yaml:"tags"`
	SortKey   string           `yaml:"-"`
	Ablative  string           `yaml:"-"`
	Letters   []*letter.Letter `yaml:"-"`
	Directory string           `yaml:"-"`
}

func New(directory string, letters []*letter.Letter) *Author {
	author := parseYaml(directory)

	author.SortKey = createSortKey(author.Name)
	author.Ablative, _ = turkishsuffix.Ablative(author.Name)
	author.Letters = letters
	author.Directory = directory

	return author
}

func createSortKey(name string) string {
	parts := strings.Split(name, " ")

	surname := removePrefix(parts[len(parts)-1])
	remaining := parts[:len(parts)-1]

	return strings.Join(append([]string{surname}, remaining...), " ")
}

func removePrefix(name string) string {
	parts := strings.Split(name, "-")

	if len(parts) < 2 {
		return name
	} else {
		return parts[len(parts)-1]
	}
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
