package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Notebook struct {
	dirpath  string
	filetype string
	notes    map[string]Note
}

func (nb *Notebook) NewID(title string) string {
	now := time.Now()
	id := now.Format("200601021504")

	return id
}

func (nb *Notebook) Add(id, title, content string) error {
	if _, exists := nb.notes[id]; exists {
		fmt.Errorf("File ID: %v already exists", id)
	}

	header := fmt.Sprintf("/*\nID:%s\nTITLE:%s\n*/", id, title)
	content = header + content

	err := os.WriteFile(filepath.Join(nb.dirpath, id+nb.filetype), []byte(content), 0644)
	if err != nil {
		return err
	}

	note := Note{
		id:      id,
		title:   title,
		content: content,
	}
	nb.notes[id] = note

	return nil
}

func (nb *Notebook) Edit(id string, newContent string) error {
	note, exists := nb.notes[id]
	if !exists {
		return fmt.Errorf("ID %v does not exist", id)
	}

	err := os.WriteFile(filepath.Join(nb.dirpath, id+nb.filetype), []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("File %s could not be edited: %v", id, err)
	}

	note.content = newContent
	nb.notes[id] = note

	return nil
}

func (nb *Notebook) Remove(id string) error {
	if _, exists := nb.notes[id]; !exists {
		return fmt.Errorf("ID %v does not exist", id)
	}

	err := os.Remove(filepath.Join(nb.dirpath, id+nb.filetype))
	if err != nil {
		return fmt.Errorf("File %s could not be deleted: %v", id, err)
	}

	delete(nb.notes, id)

	return nil
}

func (nb *Notebook) List() map[string]Note {
	return nb.notes
}

func NewNotebook(parentDir, name, filetype string) (Notebook, error) {
	if _, err := os.Stat(parentDir); err != nil {
		return Notebook{}, fmt.Errorf("Invalid Path: %v", err)
	}
	path := filepath.Join(parentDir, name)

	err := os.Mkdir(path, 0644)
	if err != nil {
		return Notebook{}, fmt.Errorf("Notebook could not be created: %v", err)
	}

	return Notebook{
		dirpath:  path,
		filetype: filetype,
		notes:    make(map[string]Note),
	}, nil
}

func LoadNotebook(path, filetype string) (Notebook, error) {
	notebookDirectory, err := os.ReadDir(path)

	if err != nil {
		return Notebook{}, fmt.Errorf("Specified Notebook could not be loaded: %v", err)
	}

	notebook := Notebook{
		dirpath:  path,
		filetype: filetype,
		notes:    make(map[string]Note),
	}

	for _, entry := range notebookDirectory {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		filepath := filepath.Join(path, name)

		file, err := os.Open(filepath)
		if err != nil {
			return Notebook{}, fmt.Errorf("File %s could not be read: %v", name, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		note := Note{
			filepath: filepath,
		}

		for scanner.Scan() {
			line := scanner.Text()

			switch true {
			case strings.HasPrefix(line, "id:"):
				note.id = cleanValue(line, "id:")

			case strings.HasPrefix(line, "title:"):
				note.title = cleanValue(line, "title")

			case line == "*/":
				break
			}
		}

		notebook.notes[note.id] = note
	}
	return notebook, nil
}

func cleanValue(line, prefix string) string {
	trimmed := strings.TrimPrefix(strings.ToLower(line), prefix)
	return strings.TrimSpace(trimmed)
}
