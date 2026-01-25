package zettel

import (
	"fmt"
	"os"
	"os/exec"
)

type Note struct {
	Id       string
	Title    string
	Content  string
	Filepath string
}

func (n *Note) OpenEditor(editor string) error {
	cmd := exec.Command(editor, n.Filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	updatedContent, err := readFile(n.Filepath)
	if err != nil {
		return fmt.Errorf("failed to sync file to notebook: %w", err)
	}

	n.Content = string(updatedContent)
	return nil
}

func readFile(filepath string) (content string, err error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	content = string(data)
	return content, nil
}
