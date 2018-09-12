package goey

import (
	"bitbucket.org/rj/goey/base"
)

var (
	tabsKind = base.NewKind("bitbucket.org/rj/goey.Tabs")
)

// Tabs describes a widget that shows a tabs.
//
// The size of the control will match the size of the currently selected child
// element, although padding will added as required to provide space for the
// border and the tabs.  However, when the user switches tabs, a relayout of
// the entire window is not forced.
//
// When calling UpdateProps, setting Value to an integer less than zero will
// leave the currently selected tab unchanged.
type Tabs struct {
	Value    int       // Index of the selected tab
	Children []TabItem // Description of the tabs

	OnChange func(int) // OnChange will be called whenever the user selects a different tab
}

// TabItem describes a tab for a Tab widget.
type TabItem struct {
	Caption string      // Text to describe the contents of this tab
	Child   base.Widget // Child widget for the tab
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Tabs) Kind() *base.Kind {
	return &tabsKind
}

// Mount creates a tabs control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Tabs) Mount(parent base.Control) (base.Element, error) {
	// Ensure that the Value is a useable index.
	w.UpdateValue()
	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateValue ensures that the index for the currently selected tab is with
// the allowed range.
func (w *Tabs) UpdateValue() {
	if w.Value >= len(w.Children) {
		w.Value = len(w.Children) - 1
	}
}

func (*tabsElement) Kind() *base.Kind {
	return &tabsKind
}

func (w *tabsElement) UpdateProps(data base.Widget) error {
	// Cast to correct type.
	tabs := data.(*Tabs)
	// Ensure that the Value is a useable index.
	tabs.UpdateValue()
	// Forward to the platform-dependant code
	return w.updateProps(tabs)
}
