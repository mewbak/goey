package dialog

import (
	"errors"
	"strings"
)

// Message is a builder to construct a message dialog to the user.
type Message struct {
	handle uintptr
	text   string
	title  string
	icon   uint
	err    error
}

// NewMessage initializes a new message object with the specified text.
// Use of the method Message on an existing Window is preferred, as the
// message can be set as a child of the top-level window.
func NewMessage(text string) *Message {
	text = strings.TrimSpace(text)
	if text == "" {
		return &Message{err: errors.New("Invalid argument, 'text' cannot be empty in call to NewMessage")}
	}
	return &Message{text: text, title: "goey"}
}

// Show completes building of the message, and shows the message to the user.
func (m *Message) Show() error {
	if m.err != nil {
		return m.err
	}

	return m.show()
}

// WithError adds an icon to the message indicating that an error has occurred.
func (m *Message) WithError() *Message {
	m.withError()
	return m
}

// WithWarn adds an icon to the message indicating that the message is a warning.
func (m *Message) WithWarn() *Message {
	m.withWarn()
	return m
}

// WithInfo adds an icon to the message indicating that the message is informational.
func (m *Message) WithInfo() *Message {
	m.withInfo()
	return m
}

// WithTitle adds a title to the message's dialog.
func (m *Message) WithTitle(text string) *Message {
	text = strings.TrimSpace(text)
	if text == "" {
		m.err = errors.New("Invalid argument, 'text' cannot be empty in call to NewMessage")
	} else {
		m.title = text
	}
	return m
}
