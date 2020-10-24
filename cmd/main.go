package main

import (
	"context"
	"github.com/ilhom9045/wallet/pkg/search"
	"log"
)

func main() {
	root := context.Background()
	ctx, cancel := context.WithCancel(root)
	files := []string{
		"data/ex.txt",
		"data/import",
		"data/import",
		"data/import",
		"data/import",
		"data/import",
		"data/import.txt",
		"data/import.txt",
		"data/export.txt",
	}
	log.Println(<-search.Any(ctx, "2;", files))
	cancel()
}
