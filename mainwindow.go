package goey

import (
	"errors"
	"fmt"
	"image"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	// ErrSetChildrenNotReentrant is returned if a reentrant call to
	// the method SetChildren is called.
	ErrSetChildrenNotReentrant = errors.New("method SetChildren is not reentrant")

	insideSetChildren uintptr
)

// Window represents a top-level window that contain other widgets.
type Window struct {
	windowImpl
}

// NewWindow create a new top-level window for the application.
func NewWindow(title string, child Widget) (*Window, error) {
	return newWindow(title, child)
}

// Close destroys the window, and releases all associated resources.
func (w *Window) Close() {
	w.close()
}

// Child returns the mounted child for the window.  In general, this
// method should not be used.
func (w *Window) Child() Element {
	return w.getChild()
}

// children assumes that the direct child of the window is a VBox, and then
// returns the children of that element.  It is used for testing.
func (w *Window) children() []Element {
	child := w.getChild()
	if child == nil {
		return nil
	}

	if vbox, ok := child.(*mountedVBox); ok {
		return vbox.children
	}

	return nil
}

func (w *windowImpl) layoutChild(windowSize Size) Size {
	// Create the constraints
	constraints := Tight(windowSize)

	// Relax maximum size when scolling is allowed
	if w.horizontalScroll {
		constraints.Max.Width = Inf
	}
	if w.verticalScroll {
		constraints.Max.Height = Inf
	}

	// Perform layout
	size := w.child.Layout(constraints)
	if !constraints.IsSatisfiedBy(size) {
		fmt.Println("constraints not satisfied,", constraints, ",", size)
	}
	return size
}

// Message returns a message constructor that can be used to build and then
// show a dialog box with a message.
func (w *Window) Message(text string) *Message {
	ret := NewMessage(text)
	w.message(ret)
	return ret
}

// Scroll returns the flags that determine whether scrolling is allowed in the
// horizontal and vertical directions.
func (w *Window) Scroll() (horizontal, vertical bool) {
	return w.horizontalScroll, w.verticalScroll
}

func (w *Window) scrollDefaults() (horizontal, vertical bool) {
	env := os.Getenv("GOEY_SCROLL")
	if env == "" {
		return false, false
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil || value >= 4 {
		return false, false
	}

	return (value & 2) == 2, (value & 1) == 1
}

// SetChild changes the child widget of the window.  As
// necessary, GUI widgets will be created or destroyed so that the GUI widgets
// match the widgets described by the parameter children.  The
// position of contained widgets will be updated to match the new layout
// properties.
func (w *Window) SetChild(child Widget) error {
	// One source of bugs in widgets is when the fire an event when being
	// updated.  This can lead to reentrant calls to SetChildren, typically
	// with incorrect information since the GUI is in an inconsistent state
	// when the event fires.  In short, this method is not reentrant.
	// The following will block changes to different windows, although
	// that shouldn't be susceptible to the same bugs.  Users in that
	// case should use Do to delay updates to other windows, but it shouldn't
	// happen in practice.
	if !atomic.CompareAndSwapUintptr(&insideSetChildren, 0, 1) {
		return ErrSetChildrenNotReentrant
	}
	defer func() {
		atomic.StoreUintptr(&insideSetChildren, 0)
	}()

	// Defer to the platform-specific code
	return w.setChild(child)
}

// SetIcon changes the icon associated with the window.
func (w *Window) SetIcon(img image.Image) error {
	return w.setIcon(img)
}

// SetScroll sets whether scrolling is allowed in the horizontal and vertical directions.
func (w *Window) SetScroll(horizontal, vertical bool) {
	w.setScroll(horizontal, vertical)
}

// SetTitle changes the caption in the title bar for the main window.
func (w *Window) SetTitle(title string) error {
	return w.setTitle(title)
}
