package main

import (
	"flag"
	"fmt"
	"strings"
)

func commandEdit(notebook Notebook, args []string) {
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	editTitlePtr := editCmd.String("nt", "", "New title for note")
	editorPtr := editCmd.String("e", "", "Editor to use")
	editCmd.Parse(args)

	editor := getEditor(*editorPtr)

	if editCmd.NArg() == 0 {
		fmt.Println("Error: pleas provide search term or ID to edit")
		return
	}
	searchTerm := strings.Join(editCmd.Args(), " ")

	query := Query{id: searchTerm, title: searchTerm, word: searchTerm}
	results := notebook.Search(query)

	var targetNote *Note

	switch {
	case len(results) == 0:
		fmt.Println("No notes found matching '%s'\n", searchTerm)
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
		targetNote.title = *editTitlePtr
	}

	err := notebook.Edit(targetNote.id, targetNote.content)
	if err != nil {
		fmt.Printf("Error saving changes to disk: %w\n", err)
		return
	}

	fmt.Printf("Note '%s' (ID: %s) updated successfully!\n", targetNote.title, targetNote.id)

}

func promptSelcection(results []*Note) (*Note, error) {
	fmt.Println("Multiple notes found. Pleas select one:")
	for i, note := range results {
		fmt.Printf("[%d] %s (ID: %s)\n", i+1, note.title, note.id)
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
