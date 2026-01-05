package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func createNotebook(name string, loc string) (Notebook, error) {
	if _, err := os.Stat(loc); err != nil {
		return Notebook{}, fmt.Errorf("Invalid Path: %v", err)
	}
	path := filepath.Join(loc, name)

	err := os.Mkdir(path, 0700)
	if err != nil {
		return Notebook{}, fmt.Errorf("Notebook could not be created: %v", err)
	}

	notebook := NewNotebook(path)

	return notebook, nil
}

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

func newID(title string) string {
	year, month, day := time.Now().Date()
	id := fmt.Sprintf("%s_%v-%v-%v", title, year, int(month), day)
	return id
}

func addNote(title string, content string, notebook Notebook) error {
	basepath := notebook.dirpath
	id := newID(title)

	err := os.WriteFile(filepath.Join(basepath, id+filetype), []byte(content), 0700)
	if err != nil {
		return err
	}
	notebook.Add(id, content)

	return nil
}

func deleteNote(id string, notebook Notebook) error {
	basepath := notebook.dirpath

	err := os.Remove(filepath.Join(basepath, id+filetype))
	if err != nil {
		return fmt.Errorf("File %s could not be deleted: %v", id, err)
	}
	return nil
}
