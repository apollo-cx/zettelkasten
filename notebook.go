package main

import (
	"fmt"
)

type Notebook struct {
	dirpath string
	notes   map[string]*Note
}

func (nb *Notebook) Add(id, content string) {
	note := Note{
		content: content,
	}
	nb.notes[id] = &note
}

func (nb *Notebook) Edit(id string, newContent string) error {
	_, exists := nb.notes[id]
	if !exists {
		return fmt.Errorf("ID %v does not exist", id)
	}
	nb.notes[id].content = newContent
	return nil
}

func (nb *Notebook) Remove(id string) error {
	_, exists := nb.notes[id]
	if !exists {
		return fmt.Errorf("ID %v does not exist", id)
	}
	delete(nb.notes, id)
	return nil
}

func (nb *Notebook) List() error {
	if len(nb.notes) == 0 {
		return fmt.Errorf("No existing notes")
	}
	for id, note := range nb.notes {
		fmt.Printf("NoteID: %v\nNote Content:\n%v", id, note.content)
		fmt.Println()
	}
	return nil
}

func NewNotebook(dir string) Notebook {
	return Notebook{
		dirpath: dir,
		notes:   make(map[string]*Note),
	}
}
