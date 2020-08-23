package main

import (
	"os"

	"peterparada.com/online-bookmarks/infrastructure"
)

func main() {
	infrastructure.Load()

	var port string = os.Getenv("ADDRESS")

	if port == "" {
		port = ":2999"
	}
}