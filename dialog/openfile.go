package dialog

import (
	"errors"
	"strings"
)

// OpenFile is a builder to construct an open file dialog to the user.
type OpenFile struct {
	handle uintptr
	title  string
	err    error
}

// NewOpenFile initializes a new open file dialog.
// Use of the method OpenFileDialog on an existing Window is preferred, as the
// message can be set as a child of the top-level window.
func NewOpenFile() *OpenFile {
	return &OpenFile{title: "goey"}
}

// Show completes building of the message, and shows the message to the user.
func (m *OpenFile) Show() (string, error) {
	if m.err != nil {
		return "", m.err
	}

	return m.show()
}

func (m *OpenFile) WithTitle(text string) *OpenFile {
	text = strings.TrimSpace(text)
	if text == "" {
		m.err = errors.New("Invalid argument, 'text' cannot be empty in call to NewMessage")
	} else {
		m.title = text
	}
	return m
}
