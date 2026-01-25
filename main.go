package main

import (
	"fmt"
	"os"

	"github.com/apollo-cx/zettelkasten/cmd"
	"github.com/apollo-cx/zettelkasten/internal/zettel"
)

const (
	notebookPath = "./notebook" // Adjust this to your preferred path
	filetype     = ".txt"
)

func main() {
	notebook, err := zettel.LoadNotebook(notebookPath, filetype)
	if err != nil {
		fmt.Printf("Error loading notebook: %v\n", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Expected 'add', 'edit', or 'search' subcommands")
		return
	}

	switch os.Args[1] {
	case "add":
		cmd.CommandAdd(notebook, os.Args[2:])
	case "edit":
		cmd.CommandEdit(notebook, os.Args[2:])
	case "search":
		cmd.CommandSearch(notebook, os.Args[2:])
	default:
		fmt.Println("Unknown command")
	}
}
