package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crgimenes/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

func main() {

	token := os.Getenv("token")
	config := dropbox.NewConfig(token)

	arqs := files.NewSearchArg("", "teste")
	nodes, more, start, err := dropbox.Search(config, arqs)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("More:", more, "\n")
	fmt.Print("Start:", start, "\n")

	for k, v := range nodes {
		fmt.Printf("%v %v\n", k, v.Name)
	}

}
