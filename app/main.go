package main

import (
	"fmt"

	"github.com/jaketreacher/gokbilgin-wiki-generator/wikiclient"
)

func main() {
	client := wikiclient.New("http://localhost:8080/api.php")

	client.Login("user", "pass")
	defer client.Logout()

	client.Edit("Test Page", "2 - My other content")
}
