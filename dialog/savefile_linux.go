package dialog

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func (m *SaveFile) show() (string, error) {
	dlg, err := gtk.FileChooserNativeDialogNew(m.title, (*gtk.Window)(unsafe.Pointer(m.handle)), gtk.FILE_CHOOSER_ACTION_SAVE, "_Save", "_Cancel")
	if err != nil {
		return "", err
	}
	defer dlg.Destroy()

	for _, v := range m.filters {
		addFilterToDialog(dlg, v.name, v.pattern)
	}

	dlg.SetFilename(m.filename)
	rc := dlg.Run()
	if gtk.ResponseType(rc) != gtk.RESPONSE_ACCEPT {
		return "", nil
	}
	return dlg.GetFilename(), nil
}

// WithParent sets the parent of the dialog box.
func (m *SaveFile) WithParent(parent *gtk.Window) *SaveFile {
	m.handle = uintptr(unsafe.Pointer(parent))
	return m
}
