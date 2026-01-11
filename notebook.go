package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	fullFileContent := header + content

	err := os.WriteFile(filepath.Join(nb.dirpath, id+nb.filetype), []byte(fullFileContent), 0644)
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

	header := fmt.Sprintf("/*\nID:%s\nTITLE:%s\n*/", id, note.title)
	fullFileContent := header + newContent

	err := os.WriteFile(filepath.Join(nb.dirpath, id+nb.filetype), []byte(fullFileContent), 0644)
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

func (nb *Notebook) Search(query Query) []Note {
	type scoredNote struct {
		score int
		note  Note
	}
	results := []scoredNote{}

	for _, note := range nb.notes {
		score := 0
		if strings.ToLower(note.id) == query.id && query.id != "" {
			score += 100
		}

		if strings.Contains(strings.ToLower(note.title), query.title) && query.title != "" {
			if len(note.title) == len(query.title) {
				score += 100
			} else {
				score += 50
			}
		}

		if strings.Contains(strings.ToLower(note.content), query.word) && query.word != "" {
			score += 10
		}

		if score > 0 {
			results = append(results,
				scoredNote{
					score: score,
					note:  note,
				},
			)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	var sortedNotes []Note
	for _, sn := range results {
		sortedNotes = append(sortedNotes, sn.note)
	}

	return sortedNotes
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

		isScanningHeader, isScanningBody := false, false
		content := ""
		for scanner.Scan() {
			line := scanner.Text()

			if isScanningBody {
				content += line + "\n"
				continue
			} else if isScanningHeader {
				switch {
				case line == "*/":
					isScanningHeader, isScanningBody = false, true
				case strings.HasPrefix(strings.ToLower(line), "id:"):
					note.id = cleanValue(line, "id:")
				case strings.HasPrefix(strings.ToLower(line), "title:"):
					note.title = cleanValue(line, "title:")
				}
			} else {
				if line == "/*" {
					isScanningHeader, isScanningBody = true, false
				}
				continue
			}
		}

		note.content = content

		if note.id == "" {
			fmt.Printf("Could'nt find note id in file: %s", note.filepath)
			continue
		}

		notebook.notes[note.id] = note
	}
	return notebook, nil
}

func cleanValue(line, prefix string) string {
	trimmed := strings.TrimSpace(line[len(prefix):])
	return trimmed
}

type Query struct {
	id    string
	title string
	word  string
}

func BuildQuery(string)
