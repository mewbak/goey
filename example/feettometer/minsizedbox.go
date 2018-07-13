package main

import (
	"bitbucket.org/rj/goey"
)

var (
	minsizedboxKind = goey.NewKind("bitbucket.org/rj/goey/example/feettometer.MinSizedBox")
)

type MinSizedBox struct {
	Child goey.Widget // Child widget.
}

// Kind returns the concrete type for use in the.Widget interface.
// Users should not need to use this method directly.
func (*MinSizedBox) Kind() *goey.Kind {
	return &minsizedboxKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *MinSizedBox) Mount(parent goey.Control) (goey.Element, error) {
	// Mount the child
	child, err := w.Child.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &minsizedboxElement{
		parent: parent,
		child:  child,
	}, nil
}

type minsizedboxElement struct {
	parent    goey.Control
	child     goey.Element
	childSize goey.Size
}

func (w *minsizedboxElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*minsizedboxElement) Kind() *goey.Kind {
	return &minsizedboxKind
}

func (w *minsizedboxElement) Layout(bc goey.Constraint) goey.Size {
	if w.child == nil {
		return bc.Constrain(goey.Size{})
	}

	width := w.child.MinIntrinsicWidth(0)
	height := w.child.MinIntrinsicHeight(width)

	size := bc.Constrain(goey.Size{width, height})
	return w.child.Layout(goey.Tight(size))
}

func (w *minsizedboxElement) MinIntrinsicHeight(width goey.Length) goey.Length {
	if w.child == nil {
		return 0
	}

	return w.child.MinIntrinsicHeight(width)
}

func (w *minsizedboxElement) MinIntrinsicWidth(height goey.Length) goey.Length {
	if w.child == nil {
		return 0
	}

	return w.child.MinIntrinsicWidth(height)
}

func (w *minsizedboxElement) SetBounds(bounds goey.Rectangle) {
	if w.child != nil {
		w.child.SetBounds(bounds)
	}
}

func (w *minsizedboxElement) updateProps(data *MinSizedBox) (err error) {
	w.child, err = goey.DiffChild(w.parent, w.child, data.Child)
	w.childSize = goey.Size{}
	return err
}

func (w *minsizedboxElement) UpdateProps(data goey.Widget) error {
	return w.updateProps(data.(*MinSizedBox))
}
