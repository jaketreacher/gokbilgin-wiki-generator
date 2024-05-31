package main

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	client.SetLanguage("tur")
	defer client.Close()
	client.SetImage("/data/input.jpg")
	text, _ := client.Text()
	fmt.Println(text)
}
