package goey

// NewMainWindow create a new top-level window for the application.
func NewMainWindow(title string, children []Widget) (*MainWindow, error) {
	return newMainWindow(title, children)
}

// Close destroys the windows, and releases all associated resources.
func (w *MainWindow) Close() {
	w.close()
}

// SetAlignment changes the vertical and horiztonal alignment properties of
// the window.
func (w *MainWindow) SetAlignment(main MainAxisAlign, cross CrossAxisAlign) {
	w.setAlignment(main, cross)
}

// SetChildren changes the child windows and widgets of the main window.
func (w *MainWindow) SetChildren(children []Widget) error {
	return w.setChildren(children)
}

func (w *MainWindow) SetTitle(title string) error {
	return w.setTitle(title)
}
