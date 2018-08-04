package main

import (
	"log"
	"os"

	"github.com/crgimenes/dropbox"
)

func main() {
	token := os.Getenv("token")
	config := dropbox.NewConfig(token)
	err := dropbox.Download(config, os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
