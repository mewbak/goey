package goey

import (
	"bitbucket.org/rj/goey/base"
)

var (
	progressKind = base.NewKind("bitbucket.org/rj/goey.Progress")
)

// Progress describes a widget that shows a progress bar.
// The model for the value is an int.
//
// If both Min and Max are zero, then Max will be updated to 100.  Other cases
// where Min == Max are not allowed.
type Progress struct {
	Value    int
	Min, Max int
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Progress) Kind() *base.Kind {
	return &progressKind
}

// Mount creates a progress control in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Progress) Mount(parent base.Control) (base.Element, error) {
	// Fill in default values for the range.
	w.UpdateRange()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

// UpdateRange sets a default range when Min and Max are uninitialized.
func (w *Progress) UpdateRange() {
	if w.Min == 0 && w.Max == 0 {
		w.Max = 100
	}
}

func (*progressElement) Kind() *base.Kind {
	return &progressKind
}

func (w *progressElement) UpdateProps(data base.Widget) error {
	pb := data.(*Progress)

	// Fill in default values for the range.
	pb.UpdateRange()
	// Forward to the platform-dependant code
	return w.updateProps(pb)
}
