package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/apollo-cx/zettelkasten/internal/zettel"
)

func commandEdit(notebook zettel.Notebook, args []string) {
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editTitlePtr := editCmd.String("nt", "", "New title for note")
	editorPtr := editCmd.String("e", "", "Editor to use")
	editCmd.Parse(args)

	editor := GetEditor(*editorPtr)

	if editCmd.NArg() == 0 {
		fmt.Println("Error: pleas provide search term or ID to edit")
		return
	}
	searchTerm := strings.Join(editCmd.Args(), " ")

	query := zettel.Query{Id: searchTerm, Title: searchTerm, Word: searchTerm}
	results := notebook.Search(query)

	var targetNote *zettel.Note

	switch {
	case len(results) == 0:
		fmt.Println("No notes found matching '%v'\n", searchTerm)
		return
	case len(results) == 1:
		targetNote = results[0]
	default:
		var err error
		targetNote, err = promptSelcection(results)
		if err != nil {
			fmt.Printf("Selection error: %v\n", err)
			return
		}
	}

	if err := targetNote.OpenEditor(editor); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if *editTitlePtr != "" {
		targetNote.Title = *editTitlePtr
	}

	err := notebook.Edit(targetNote.Id, targetNote.Content)
	if err != nil {
		fmt.Printf("Error saving changes to disk: %w\n", err)
		return
	}

	fmt.Printf("Note '%s' (ID: %s) updated successfully!\n", targetNote.Title, targetNote.Id)

}

func promptSelcection(results []*zettel.Note) (*zettel.Note, error) {
	fmt.Println("Multiple notes found. Pleas select one:")
	for i, note := range results {
		fmt.Printf("[%d] %s (ID: %s)\n", i+1, note.Title, note.Id)
	}

	fmt.Print("Chose by entering a number: ")
	var choice int
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		return nil, fmt.Errorf("invalid input")
	}

	if choice < 1 || choice > len(results) {
		return nil, fmt.Errorf("choice out of range")
	}

	return results[choice-1], nil
}
