package main

import (
	"flag"
	"fmt"
)

func commandAdd(notebook Notebook, args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	titlePtr := addCmd.String("t", "", "Title of note")
	editorPtr := addCmd.String("e", "", "Editor to use")
	addCmd.Parse(args)

	editor := getEditor(*editorPtr)

	id := notebook.NewID()
	err := notebook.Add(*titlePtr, "")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}

	note, ok := notebook.notes[id]
	if !ok {
		fmt.Printf("Error: Note with ID %s not found in notebook\n", id)
		return
	} else if err := note.OpenEditor(editor); err != nil {
		fmt.Printf("Error during initial edit: %v\n", err)
		return
	}

	fmt.Printf("Note '%s' created and saved successfully!", note.title)
}
