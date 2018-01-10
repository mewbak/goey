package goey

import (
	"errors"
	"strings"
)

type Message struct {
	handle  uintptr
	text    string
	title string
	icon    uint
	err     error
}

func NewMessage(text string) *Message {
	text = strings.TrimSpace(text)
	if text == "" {
		return &Message{err: errors.New("Invalid argument, 'text' cannot be empty in call to NewMessage")}
	}
	return &Message{text: text, title: "goey"}
}
func (m *Message) Show() error {
	if m.err != nil {
		return m.err
	}

	return m.show()
}

func (m *Message) WithError() *Message {
	m.withError()
	return m
}

func (m *Message) WithWarn() *Message {
	m.withWarn()
	return m
}

func (m *Message) WithInfo() *Message {
	m.withInfo()
	return m
}

func (m *Message) WithTitle(text string) *Message {
	text = strings.TrimSpace(text)
	if text == "" {
		m.err= errors.New("Invalid argument, 'text' cannot be empty in call to NewMessage")
	} else {
		m.title = text
	}
	return m
}
