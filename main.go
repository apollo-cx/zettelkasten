package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	filetype = ".txt"
)

func main() {

}

/* Notes Stuff */

type Notebook struct {
	dirpath string
	notes   map[string]*Note
}

type Note struct {
	content string
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

func NewNotebook(dirpath string) Notebook {
	return Notebook{
		dirpath: dirpath,
		notes:   make(map[string]*Note),
	}
}

func NewFileID(title string) string {
	year, month, day := time.Now().Date()
	id := fmt.Sprintf("%s_%v-%v-%v", title, year, int(month), day)
	return id
}

/* Persictance stuff */

func loadNotebook(dir string) (Notebook, error) {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return Notebook{}, fmt.Errorf("Specified Notebook could not be loaded: %v", err)
	}

	notebook := NewNotebook(dir)

	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return Notebook{}, fmt.Errorf("File %s could not be read: %v", name, err)
		}

		notebook.Add(name, string(data))
	}
	return notebook, nil
}

func addNote(title string, content string, notebook Notebook) error {
	basepath := notebook.dirpath
	id := NewFileID(title)

	err := os.WriteFile(filepath.Join(basepath, id+filetype), []byte(content), 0700)
	if err != nil {
		return err
	}
	notebook.Add(id, content)

	return nil
}
