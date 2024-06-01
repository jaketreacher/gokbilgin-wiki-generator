package main

import (
	"fmt"

	"github.com/jaketreacher/gokbilgin-wiki-generator/wikiclient"
)

func main() {
	client := wikiclient.New("http://localhost:8080/api.php")

	token := client.TokenQuery()

	fmt.Println(token)
}
