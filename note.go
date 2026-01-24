package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Note struct {
	id       string
	title    string
	content  string
	filepath string
}

func (n *Note) OpenEditor(editor string) error {
	cmd := exec.Command(editor, n.filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	updatedContent, err := readFile(n.filepath)
	if err != nil {
		return fmt.Errorf("failed to sync file to notebook: %w", err)
	}

	n.content = string(updatedContent)
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
