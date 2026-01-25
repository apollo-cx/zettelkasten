package cmd

import (
	"flag"
	"fmt"

	"github.com/apollo-cx/zettelkasten/internal/zettel"
)

func CommandAdd(notebook zettel.Notebook, args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	titlePtr := addCmd.String("t", "", "Title of note")
	editorPtr := addCmd.String("e", "", "Editor to use")
	addCmd.Parse(args)

	editor := GetEditor(*editorPtr)

	note, err := notebook.Add(*titlePtr, "")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	if err := note.OpenEditor(editor); err != nil {
		fmt.Printf("Error during initial edit: %v\n", err)
		return
	}

	fmt.Printf("Note '%s' created and saved successfully!", note.Title)
}
