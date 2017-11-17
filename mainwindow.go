package goey

// NewMainWindow create a new top-level window for the application.
func NewMainWindow(title string, children []Widget) (*MainWindow, error) {
	return newMainWindow(title, children)
}
