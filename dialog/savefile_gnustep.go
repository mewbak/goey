// +build gnustep

package dialog

import (
	"bitbucket.org/rj/goey/cocoa"
)

func (m *SaveFile) show() (string, error) {
	cocoa.SavePanel(m.parent)
	return "", nil
}

// WithParent sets the parent of the dialog box.
func (m *SaveFile) WithParent(parent *cocoa.Window) *SaveFile {
	m.parent = parent
	return m
}
