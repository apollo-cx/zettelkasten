package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/apollo-cx/zettelkasten/internal/zettel"
)

func commandSearch(notebook zettel.Notebook, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: nothing to search for")
		return
	}
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	idPtr := searchCmd.String("i", "", "ID to search for")
	titlePtr := searchCmd.String("t", "", "Title to search for")
	wordPtr := searchCmd.String("w", "", "Word to search for")

	searchCmd.Parse(args)
	query := zettel.Query{}

	if *idPtr != "" {
		query.Id = *idPtr
	}
	if *titlePtr != "" {
		query.Title = *titlePtr
	}
	if *wordPtr != "" {
		query.Word = *wordPtr
	}

	if searchCmd.NArg() > 0 {
		// Join the remaining words ("mars", "rover") into "mars rover"
		term := strings.Join(searchCmd.Args(), " ")

		// If no specific flags were set, search for this term everywhere
		if query.Id == "" && query.Title == "" && query.Word == "" {
			query.Id = term
			query.Title = term
			query.Word = term
		}
	}

	results := notebook.Search(query)

	fmt.Printf("%d results:\n", len(results))
	for _, result := range results {
		fmt.Printf("%s -> %s\n", result.Title, result.Filepath)
	}
}
