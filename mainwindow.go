package goey

// Window represents a top-level window that contain other widgets.
type Window struct {
	windowImpl
}

// NewWindow create a new top-level window for the application.
func NewWindow(title string, children []Widget) (*Window, error) {
	return newWindow(title, children)
}

// Close destroys the window, and releases all associated resources.
func (w *Window) Close() {
	w.close()
}

// SetAlignment changes the vertical and horiztonal alignment properties of
// the window.  These properties affect the layout of child widgets.  The
// main axis for alignment is vertical, with the cross axis being horizontal.
func (w *Window) SetAlignment(main MainAxisAlign, cross CrossAxisAlign) {
	w.setAlignment(main, cross)
}

// SetChildren changes the child windows and widgets of the main window.
func (w *Window) SetChildren(children []Widget) error {
	return w.setChildren(children)
}

// SetTitle changes the caption in the title bar for the main window.
func (w *Window) SetTitle(title string) error {
	return w.setTitle(title)
}
