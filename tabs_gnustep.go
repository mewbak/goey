// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type tabsElement struct {
	control *cocoa.TabView
	value   int
	child   base.Element
	widgets []TabItem
	insets  Insets

	cachedBounds base.Rectangle
	cachedTabsW  base.Length
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTabView(parent.Handle)
	for _, v := range w.Children {
		control.AddItem(v.Caption)
	}

	child := base.Element(nil)
	if len(w.Children) > 0 {
		parent := base.Control{&control.View}
		child_, err := base.Mount(parent, w.Children[w.Value].Child)
		if err != nil {
			control.Close()
			return nil, err
		}
		child = child_
		control.SelectItem(w.Value)
	}

	retval := &tabsElement{
		control: control,
		child:   child,
		value:   w.Value,
		widgets: w.Children,
	}
	retval.setInsets(w.Insets)
	return retval, nil
}

func (w *tabsElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		w.cachedTabsW = 100 * base.DIP
	}
	return w.cachedTabsW
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())

	if w.child != nil {
		// Determine the bounds for the child widget
		dx, dy := bounds.Dx(), bounds.Dy()

		bounds.Min.X = w.insets.Left
		bounds.Min.Y = w.insets.Top
		bounds.Max.X = dx - w.insets.Right
		bounds.Max.Y = dy - w.insets.Bottom

		// Update bounds for the child
		w.cachedBounds = bounds
		w.child.SetBounds(bounds)
	}
}

func (w *tabsElement) setInsets(value Insets) {
	// Adjust the insets to adjust for the content area
	x, y := w.control.ContentOrigin()
	wid, h := w.control.ContentInsets()

	value.Left += base.FromPixelsX(x)
	value.Top += base.FromPixelsY(y)
	value.Right += base.FromPixelsX(wid - x)
	value.Bottom += base.FromPixelsY(h - y)
	w.insets = value
}

func (w *tabsElement) updateProps(data *Tabs) error {
	return nil
}
