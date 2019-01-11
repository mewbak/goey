// +build gnustep

package dialog

import (
	"bitbucket.org/rj/goey/cocoa"
)

func (m *OpenFile) show() (string, error) {
	cocoa.OpenPanel(m.parent)
	return "", nil
}

// WithParent sets the parent of the dialog box.
func (m *OpenFile) WithParent(parent *cocoa.Window) *OpenFile {
	m.parent = parent
	return m
}
