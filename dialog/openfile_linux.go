package dialog

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func (m *OpenFile) show() (string, error) {
	dlg, err := gtk.FileChooserNativeDialogNew(m.title, (*gtk.Window)(unsafe.Pointer(m.handle)), gtk.FILE_CHOOSER_ACTION_OPEN, "Open", "Cancel")
	if err != nil {
		return "", err
	}
	defer dlg.Destroy()

	rc := dlg.Run()
	if gtk.ResponseType(rc) != gtk.RESPONSE_ACCEPT {
		return "", nil
	}
	return dlg.GetFilename(), nil
}

// WithParent sets the parent of the dialog box.
func (m *OpenFile) WithParent(parent *gtk.Window) *OpenFile {
	m.handle = uintptr(unsafe.Pointer(parent))
	return m
}
