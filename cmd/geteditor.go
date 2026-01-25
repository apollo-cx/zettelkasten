package cmd

import "os"

func GetEditor(flagValue string) string {
	editor := os.Getenv("EDITOR")
	if flagValue != "" {
		editor = flagValue
	} else if editor == "" {
		editor = "nano"
	}
	return editor
}
