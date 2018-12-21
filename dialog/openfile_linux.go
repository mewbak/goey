package dialog

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *OpenFile) show() (string, error) {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(m.title, m.parent, gtk.FILE_CHOOSER_ACTION_OPEN, "_Open", gtk.RESPONSE_ACCEPT, "_Cancel", gtk.RESPONSE_CANCEL)
	if err != nil {
		return "", err
	}
	activeDialogForTesting = &dlg.Dialog
	defer func() {
		activeDialogForTesting = nil
		dlg.Destroy()
	}()

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

func addFilterToDialog(dlg *gtk.FileChooserDialog, name, pattern string) error {
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
	m.parent = parent
	return m
}
