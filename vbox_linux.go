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

	halign := w.AlignCross.HAlign()
	children := make([]MountedWidget, 0, len(w.Children))
	for _, v := range w.Children {
		mountedWidget, err := v.Mount(NativeWidget{&control.Widget})
		if err != nil {
			return nil, err
		}
		vbox_crossAlign(mountedWidget, halign)
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

func vbox_crossAlign(widget MountedWidget, align gtk.Align) {
	// Label's don't stretch, and we don't want them to be centered.
	if align == gtk.ALIGN_FILL {
		if ptr, ok := widget.(*mountedLabel); ok {
			ptr.handle.SetHAlign(gtk.ALIGN_START)
			return
		}
	}
	widget.Handle().SetHAlign(align)
}

func (w *mountedVBox) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	// Save main axis alignment
	w.alignMain = main

	// Save cross axis alignment, update children
	if w.alignCross != cross {
		w.alignCross = cross
		halign := cross.HAlign()
		for _, v := range w.children {
			vbox_crossAlign(v, halign)
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
