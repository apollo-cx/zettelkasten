package zettel

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
	Dirpath  string
	Filetype string
	Notes    map[string]*Note
}

func newId() string {
	now := time.Now()
	id := now.Format("20060102150405")

	return id
}

func (nb *Notebook) Add(title, content string) (Note, error) {
	id := newId()
	if _, exists := nb.Notes[id]; exists {
		return Note{}, fmt.Errorf("File Id: %v already exists", id)
	}

	if title == "" {
		title = id
	}

	header := fmt.Sprintf("/*\nId:%s\nTitle:%s\n*/", id, title)
	fullFileContent := header + content
	filepath := filepath.Join(nb.Dirpath, id+nb.Filetype)

	err := os.WriteFile(filepath, []byte(fullFileContent), 0644)
	if err != nil {
		return Note{}, err
	}

	note := &Note{
		Id:       id,
		Title:    title,
		Content:  content,
		Filepath: filepath,
	}

	nb.Notes[id] = note
	return *note, err
}

func (nb *Notebook) Edit(id, newContent string) error {
	note, exists := nb.Notes[id]
	if !exists {
		return fmt.Errorf("Id %v does not exist", id)
	}

	header := fmt.Sprintf("/*\nId:%s\nTitle:%s\n*/", id, note.Title)
	fullFileContent := header + newContent

	err := os.WriteFile(filepath.Join(nb.Dirpath, id+nb.Filetype), []byte(fullFileContent), 0644)
	if err != nil {
		return fmt.Errorf("File %s could not be edited: %v", id, err)
	}

	note.Content = newContent
	return nil
}

func (nb *Notebook) Remove(id string) error {
	if _, exists := nb.Notes[id]; !exists {
		return fmt.Errorf("Id %v does not exist", id)
	}

	err := os.Remove(filepath.Join(nb.Dirpath, id+nb.Filetype))
	if err != nil {
		return fmt.Errorf("File %s could not be deleted: %v", id, err)
	}

	delete(nb.Notes, id)
	return nil
}

type Query struct {
	Id    string
	Title string
	Word  string
}

func (nb *Notebook) Search(query Query) []*Note {
	type scoredNote struct {
		score int
		note  *Note
	}
	results := []scoredNote{}

	for _, note := range nb.Notes {
		score := 0
		if strings.ToLower(note.Id) == query.Id && query.Id != "" {
			score += 100
		}

		if strings.Contains(strings.ToLower(note.Title), query.Title) && query.Title != "" {
			if len(note.Title) == len(query.Title) {
				score += 100
			} else {
				score += 50
			}
		}

		if strings.Contains(strings.ToLower(note.Content), query.Word) && query.Word != "" {
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

	var sortedNotes []*Note
	for _, sn := range results {
		sortedNotes = append(sortedNotes, sn.note)
	}

	return sortedNotes
}

func (nb *Notebook) List() map[string]*Note {
	return nb.Notes
}

func NewNotebook(parentDir, name, filetype string) (Notebook, error) {
	if _, err := os.Stat(parentDir); err != nil {
		return Notebook{}, fmt.Errorf("InvalId Path: %v", err)
	}
	path := filepath.Join(parentDir, name)

	err := os.Mkdir(path, 0644)
	if err != nil {
		return Notebook{}, fmt.Errorf("Notebook could not be created: %v", err)
	}

	return Notebook{
		Dirpath:  path,
		Filetype: filetype,
		Notes:    make(map[string]*Note),
	}, nil
}

func LoadNotebook(path, filetype string) (Notebook, error) {
	notebookDirectory, err := os.ReadDir(path)

	if err != nil {
		return Notebook{}, fmt.Errorf("Specified Notebook could not be loaded: %v", err)
	}

	notebook := Notebook{
		Dirpath:  path,
		Filetype: filetype,
		Notes:    make(map[string]*Note),
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
		note := &Note{
			Filepath: filepath,
		}

		isScanningHeader, isScanningBody := false, false
		Content := ""
		for scanner.Scan() {
			line := scanner.Text()

			if isScanningBody {
				Content += line + "\n"
				continue
			} else if isScanningHeader {
				switch {
				case line == "*/":
					isScanningHeader, isScanningBody = false, true
				case strings.HasPrefix(strings.ToLower(line), "Id:"):
					note.Id = cleanValue(line, "Id:")
				case strings.HasPrefix(strings.ToLower(line), "Title:"):
					note.Title = cleanValue(line, "Title:")
				}
			} else {
				if line == "/*" {
					isScanningHeader, isScanningBody = true, false
				}
				continue
			}
		}

		note.Content = Content

		if note.Id == "" {
			fmt.Printf("Could'nt find note Id in file: %s", note.Filepath)
			continue
		}

		notebook.Notes[note.Id] = note
	}
	return notebook, nil
}

func cleanValue(line, prefix string) string {
	trimmed := strings.TrimSpace(line[len(prefix):])
	return trimmed
}
