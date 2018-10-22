// +build !gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
)

var (
	selectKind = base.NewKind("bitbucket.org/rj/goey.SelectInput")
)

// SelectInput describes a widget that users can click to select one from a fixed list of choices.
type SelectInput struct {
	Items    []string        // Items is an array of strings representing the user's possible choices
	Value    int             // Value is the index of the currently selected item
	Unset    bool            // Unset is a flag indicating that no choice has yet been made
	Disabled bool            // Disabled is a flag indicating that the user cannot interact with this field
	OnChange func(value int) // OnChange will be called whenever the user changes the value for this field
	OnFocus  func()          // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur   func()          // OnBlur will be called whenever the field loses the keyboard focus
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*SelectInput) Kind() *base.Kind {
	return &selectKind
}

// Mount creates a select control (combobox) in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *SelectInput) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*selectinputElement) Kind() *base.Kind {
	return &selectKind
}

func (w *selectinputElement) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*SelectInput))
}
