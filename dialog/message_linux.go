package dialog

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *Message) show() error {
	dlg := gtk.MessageDialogNew(m.parent, gtk.DIALOG_MODAL, gtk.MessageType(m.icon), gtk.BUTTONS_OK, m.text)
	activeDialogForTesting = &dlg.Dialog
	defer func() {
		activeDialogForTesting = nil
		dlg.Destroy()
	}()

	dlg.SetTitle(m.title)
	dlg.Run()
	return nil
}

func (m *Message) withError() {
	m.icon = uint(gtk.MESSAGE_ERROR)
}

func (m *Message) withWarn() {
	m.icon = uint(gtk.MESSAGE_WARNING)
}

func (m *Message) withInfo() {
	m.icon = uint(gtk.MESSAGE_INFO)
}

// WithParent sets the parent of the dialog box.
func (m *Message) WithParent(parent *gtk.Window) *Message {
	m.parent = parent
	return m
}
