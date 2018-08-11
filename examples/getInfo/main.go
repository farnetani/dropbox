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

	arqs := files.NewGetMetadataArg("/teste1/arquivo.txt")
	node, err := dropbox.GetFileInfo(config, arqs)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(node.Name)

}
