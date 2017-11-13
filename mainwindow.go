package goey

// New creates and opens a new webview window using the given settings. The
// returned object implements the WebView interface. This function returns nil
// if a window can not be created.
func NewMainWindow(title string, children []Widget) (*MainWindow, error) {
	return newMainWindow(title, children)
}
