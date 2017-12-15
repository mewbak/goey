package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedVBox struct {
	NativeWidget
	children   []MountedWidget
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign
}

func (w *VBox) mount(parent NativeWidget) (MountedWidget, error) {
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
		mountedWidget.Handle().SetHAlign(w.AlignCross.HAlign())
		children = append(children, mountedWidget)
	}

	retval := &mountedVBox{
		NativeWidget: NativeWidget{&control.Widget},
		children:     children,
		alignMain:    w.AlignMain,
		alignCross:   w.AlignCross,
	}

	control.Connect("destroy", vbox_onDestroy, retval)
	control.Show()

	return retval, nil
}

func (a CrossAxisAlign) HAlign() gtk.Align {
	switch a {
	case Stretch:
		return gtk.ALIGN_FILL
	case CrossStart:
		return gtk.ALIGN_START
	case CrossCenter:
		return gtk.ALIGN_CENTER
	case CrossEnd:
		return gtk.ALIGN_END
	}

	panic("not reachable")
}

func vbox_onDestroy(widget *gtk.Box, mounted *mountedVBox) {
	mounted.handle = nil
}

func (w *mountedVBox) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	// Save main axis alignment
	w.alignMain = main

	// Save cross axis alignment, update children
	if w.alignCross != cross {
		w.alignCross = cross
		for _, v := range w.children {
			v.Handle().SetHAlign(cross.HAlign())
		}
	}

	return nil
}

func (w *mountedVBox) setChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.NativeWidget, w.children, children)
	return err
}

func (w *mountedVBox) updateProps(data *VBox) error {
	// Save new alignment
	w.setAlignment(data.AlignMain, data.AlignCross)

	// Set children
	return w.setChildren(data.Children)
}
