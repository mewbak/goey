package dialog

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

func (m *OpenFile) show() (string, error) {
	dlg, err := gtk.FileChooserNativeDialogNew(m.title, (*gtk.Window)(unsafe.Pointer(m.handle)), gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", "_Cancel")
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

func addFilterToDialog(dlg *gtk.FileChooserNativeDialog, name, pattern string) error {
	filter, err := gtk.FileFilterNew()
	if err != nil {
		return err
	}
	defer filter.Unref()

	filter.SetName(name)
	filter.AddPattern(pattern)
	dlg.AddFilter(filter)
	return nil
}

// WithParent sets the parent of the dialog box.
func (m *OpenFile) WithParent(parent *gtk.Window) *OpenFile {
	m.handle = uintptr(unsafe.Pointer(parent))
	return m
}
