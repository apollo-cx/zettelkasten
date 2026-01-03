package main

import (
	"fmt"
)

/*const (
	basePath = filepath
)*/

func main() {
	notebook := Notebook{
		notes: make(map[int]*Note),
	}

	fmt.Println("-- 1. Try addNote() --")
	notebook.addNote("hello")
	notebook.addNote("my god")
	notebook.addNote("so many notes")

	notebook.listNotes()

	fmt.Println("-- 2. Try removeNote() --")
	notebook.removeNote(1)
	notebook.removeNote(2)

	notebook.listNotes()

	fmt.Println("-- 3. Try editNote() --")
	notebook.editNote(3, "hello")

	notebook.listNotes()
}

type Notebook struct {
	notes map[int]*Note
}

type Note struct {
	content string
	path    string
}

func (nb *Notebook) addNote(content string) {
	note := Note{
		content: content,
	}
	nb.notes[len(nb.notes)+1] = &note
}

func (nb *Notebook) editNote(id int, newContent string) error {
	_, exists := nb.notes[id]
	if !exists {
		return fmt.Errorf("ID %d does not exist", id)
	}
	nb.notes[id].content = newContent
	return nil
}

func (nb *Notebook) removeNote(id int) error {
	_, exists := nb.notes[id]
	if !exists {
		return fmt.Errorf("ID %d does not exist", id)
	}
	delete(nb.notes, id)
	return nil
}

func (nb *Notebook) listNotes() error {
	if len(nb.notes) == 0 {
		return fmt.Errorf("No existing notes")
	}
	for id, note := range nb.notes {
		fmt.Printf("NoteID: %d\nNote Content:\n%s", id, note.content)
		fmt.Println()
	}
	return nil
}
