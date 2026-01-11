package main

import (
	"flag"
	"fmt"
	"strings"
)

func commandHelp() error {
	fmt.Println()
	fmt.Println("Welcom to Zettelkasten!")
	fmt.Println("Usage")
	fmt.Println()

	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	fmt.Println()
	return nil
}

func commandSearch(notebook Notebook, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: nothing to search for")
		return
	}
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	idPtr := searchCmd.String("id", "", "ID to search for")
	titlePtr := searchCmd.String("title", "", "Title to search for")
	wordPtr := searchCmd.String("word", "", "Word to search for")

	searchCmd.Parse(args)
	query := Query{}

	if *idPtr != "" {
		query.id = *idPtr
	}
	if *titlePtr != "" {
		query.title = *titlePtr
	}
	if *wordPtr != "" {
		query.word = *wordPtr
	}

	if searchCmd.NArg() > 0 {
		// Join the remaining words ("mars", "rover") into "mars rover"
		term := strings.Join(searchCmd.Args(), " ")

		// If no specific flags were set, search for this term everywhere
		if query.id == "" && query.title == "" && query.word == "" {
			query.id = term
			query.title = term
			query.word = term
		}
	}

	results := notebook.Search(query)

	fmt.Printf("%d results:\n", len(results))
	for _, result := range results {
		fmt.Printf("%s -> %s\n", result.title, result.filepath)
	}
}
