package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedHBox struct {
	NativeWidget
	children []MountedWidget
}

func (w *HBox) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 11)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	children := make([]MountedWidget, 0, len(w.Children))
	for _, v := range w.Children {
		mountedWidget, err := v.Mount(NativeWidget{&control.Container.Widget})
		if err != nil {
			return nil, err
		}
		children = append(children, mountedWidget)
	}

	retval := &mountedHBox{
		NativeWidget: NativeWidget{&control.Container.Widget},
		children:     children,
	}

	control.Connect("destroy", hbox_onDestroy, retval)
	control.Show()

	return retval, nil
}

func hbox_onDestroy(widget *gtk.Box, mounted *mountedHBox) {
	mounted.handle = nil
}

func (w *mountedHBox) SetChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.NativeWidget, w.children, children)
	return err
}
