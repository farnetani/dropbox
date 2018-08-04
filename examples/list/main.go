package main

import (
	"fmt"
	"log"
	"os"

	"github.com/crgimenes/dropbox"
)

func main() {
	token := os.Getenv("token")
	config := dropbox.NewConfig(token)
	nodes, err := dropbox.List(config, "")
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range nodes {
		fmt.Printf("%v %v\n", k, v.Name)
	}
}
