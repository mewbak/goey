package goey

import "github.com/gotk3/gotk3/gtk"

type NativeWidget struct {
	handle *gtk.Widget
}

func (w *NativeWidget) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

type NativeMountedWidget interface {
}
