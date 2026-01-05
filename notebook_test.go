package main

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	notebook := Notebook{
		notes: make(map[string]*Note),
	}
	cases := []struct {
		inputs   []string
		expected Note
	}{
		{
			inputs:   []string{"hello_26-1-5", "hello world!"},
			expected: Note{content: "hello world!"},
		},
		{
			inputs:   []string{"todo_25-12-11", "TODO:\n- do laundry\n-make dishes"},
			expected: Note{content: "TODO:\n- do laundry\n-make dishes"},
		},
		{
			inputs:   []string{"notes_25-12-11", "  --- int stands for integer ---  "},
			expected: Note{content: "  --- int stands for integer ---  "},
		},
	}

	for _, c := range cases {
		notebook.Add(c.inputs[0], c.inputs[1])

		note, exists := notebook.notes[c.inputs[0]]
		if !exists {
			t.Errorf("Created Note does not exist in %v", notebook.notes)
		} else if note.content != c.expected.content {
			t.Errorf("Contents do not match:\n %s (input)\n %s (expected)\n", note.content, c.expected.content)
		}
	}

}

func TestEdit(t *testing.T) {
	notebook := Notebook{
		notes: make(map[string]*Note),
	}

	// notes to edit
	notebook.Add("hello_11-11-11", "hello world")
	notebook.Add("todo_25-11-20", "TODO:\n- do laundry\n-make dishes")

	cases := []struct {
		inputs   []string
		expected Note
	}{
		{
			inputs:   []string{"hello_11-11-11", "hello world, i dislike spagetti"},
			expected: Note{content: "hello world, i dislike spagetti"},
		},
		{
			inputs:   []string{"todo_25-11-20", "TODO:\n- do laundry"},
			expected: Note{content: "TODO:\n- do laundry"},
		},
	}

	for _, c := range cases {
		notebook.Edit(c.inputs[0], c.inputs[1])

		note, exists := notebook.notes[c.inputs[0]]
		if !exists {
			t.Errorf("Edited Note does not exist anymore %v", notebook.notes)
		} else if note.content != c.expected.content {
			t.Errorf("Edited content does not match with expected content:\n%s (input)\n %s (expected)\n", note.content, c.expected.content)
		}
	}
}

func TestRemove(t *testing.T) {
	notebook := Notebook{
		notes: make(map[string]*Note),
	}

	// notes to edit
	notebook.Add("hello_11-11-11", "hello world")
	notebook.Add("todo_25-11-20", "TODO:\n- do laundry\n-make dishes")

	cases := []struct {
		inputs   []string
		expected bool
	}{
		{
			inputs:   []string{"hello_11-11-11", "hello world, i dislike spagetti"},
			expected: false,
		},
		{
			inputs:   []string{"todo_25-11-20", "TODO:\n- do laundry"},
			expected: false,
		},
	}

	for _, c := range cases {
		notebook.Remove(c.inputs[0])

		_, exists := notebook.notes[c.inputs[0]]
		if exists != c.expected {
			t.Errorf("Removed Note still exists: %v", notebook.notes)
		}
	}
}

func TestNewNoteBook(t *testing.T) {
	notebook := NewNotebook("/notebook")
	expected := Notebook{
		dirpath: "/notebook",
		notes:   make(map[string]*Note),
	}
	if eq := reflect.DeepEqual(notebook, expected); !eq {
		t.Errorf("Created Notebook does not match expected Notebook:\n%v (created)\n%v (expected)", notebook, expected)
	}
}
