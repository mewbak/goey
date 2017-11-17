package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type MountedVBox struct {
	NativeWidget
	children []MountedWidget
}

func (w *VBox) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 11)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	children := make([]MountedWidget, 0, len(w.Children))
	for _, v := range w.Children {
		mountedWidget, err := v.Mount(NativeWidget{&control.Widget})
		if err != nil {
			return nil, err
		}
		children = append(children, mountedWidget)
	}

	retval := &MountedVBox{
		NativeWidget: NativeWidget{&control.Widget},
		children:     children,
	}

	control.Connect("destroy", vbox_onDestroy, retval)
	control.Show()

	return retval, nil
}

func vbox_onDestroy(widget *gtk.Box, mounted *MountedVBox) {
	mounted.handle = nil
}

func (w *MountedVBox) SetChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.NativeWidget, w.children, children)
	return err
}
