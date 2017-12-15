package goey

// MainWindow holds data for a top-level window.
type MainWindow struct {
	mainWindow // Platform-specific implementation
}

// NewMainWindow create a new top-level window for the application.
func NewMainWindow(title string, children []Widget) (*MainWindow, error) {
	return newMainWindow(title, children)
}

// Close destroys the window, and releases all associated resources.
func (w *MainWindow) Close() {
	w.close()
}

// SetAlignment changes the vertical and horiztonal alignment properties of
// the window.  If necessary, the position of contained widgets will be updated
// to match new layout properties.
func (w *MainWindow) SetAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	return w.setAlignment(main, cross)
}

// SetChildren changes the child windows and widgets of the main window.  As
// necessary GUI widgets will be created or destroyed so that the GUI widgets
// match the widgets described by the parameter children.  If necessary, the
// position of contained widgets will be updated to match the new layout
// properties.
func (w *MainWindow) SetChildren(children []Widget) error {
	return w.setChildren(children)
}

// SetTitle changes the caption or title of the window.
func (w *MainWindow) SetTitle(title string) error {
	return w.setTitle(title)
}
