package dialog

import (
	"github.com/gotk3/gotk3/gtk"
)

func (m *SaveFile) show() (string, error) {
	dlg, err := gtk.FileChooserDialogNewWith2Buttons(m.title, m.parent, gtk.FILE_CHOOSER_ACTION_SAVE, "_Save", gtk.RESPONSE_ACCEPT, "_Cancel", gtk.RESPONSE_CANCEL)
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

// WithParent sets the parent of the dialog box.
func (m *SaveFile) WithParent(parent *gtk.Window) *SaveFile {
	m.parent = parent
	return m
}
