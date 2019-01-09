// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/cocoa"
	"unsafe"
)

func (m *Message) show() error {
	cocoa.MessageDialog((*cocoa.Window)(unsafe.Pointer(m.handle)), m.text, m.title, byte(m.icon))
	return nil
}

func (m *Message) withError() {
	m.icon = 'e'
}

func (m *Message) withWarn() {
	m.icon = 'w'
}

func (m *Message) withInfo() {
	m.icon = 'i'
}
