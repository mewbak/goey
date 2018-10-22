// +build !gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/gtk"
)

type hrElement struct {
	Control
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	control, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(control)

	retval := &hrElement{
		Control: Control{&control.Widget},
	}

	control.Connect("destroy", hrOnDestroy, retval)
	control.Show()

	return retval, nil
}

func hrOnDestroy(widget *gtk.Separator, mounted *hrElement) {
	mounted.handle = nil
}
