// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type tabsElement struct {
	control  *cocoa.TabView
	value    int
	child    base.Element
	widgets  []TabItem
	insets   Insets
	onChange func(int)

	cachedBounds base.Rectangle
	cachedInsets base.Size
	cachedTabsW  base.Length
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTabView(parent.Handle)
	for _, v := range w.Children {
		control.AddItem(v.Caption)
	}

	child := base.Element(nil)
	if len(w.Children) > 0 {
		parent := base.Control{control.ContentView(w.Value)}
		child_, err := base.Mount(parent, w.Children[w.Value].Child)
		if err != nil {
			control.Close()
			return nil, err
		}
		child = child_
		control.SelectItem(w.Value)
	}

	retval := &tabsElement{
		control:  control,
		child:    child,
		value:    w.Value,
		widgets:  w.Children,
		insets:   w.Insets,
		onChange: w.OnChange,
	}
	control.SetOnChange(retval.OnChange)
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

func (w *tabsElement) contentInsets() (base.Length, base.Length) {
	if w.cachedInsets.Height == 0 {
		x, y := w.control.ContentInsets()
		w.cachedInsets.Width = base.FromPixelsX(x)
		w.cachedInsets.Height = base.FromPixelsY(y)
	}

	return w.cachedInsets.Width, w.cachedInsets.Height
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		w.cachedTabsW = 100 * base.DIP
	}
	return w.cachedTabsW
}

func (w *tabsElement) mountPage(page int) {
	if w.child != nil {
		println("close")
		w.child.Close()
		w.child = nil
	}

	println("mount")
	parent := base.Control{w.control.ContentView(page)}
	child, _ := base.Mount(parent, w.widgets[page].Child)
	child.Layout(base.Tight(base.Size{
		w.cachedBounds.Dx(),
		w.cachedBounds.Dy(),
	}))
	child.SetBounds(w.cachedBounds)
	w.child = child
}

func (w *tabsElement) OnChange(page int) {
	println("...", page)
	if page != w.value {
		if w.onChange != nil {
			println("cb")
			w.onChange(page)
		}
		if page != w.value {
			println("mount")
			w.mountPage(page)
			w.value = page
		}
	}
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())

	if w.child != nil {
		// Determine the bounds for the child widget
		dx := bounds.Dx() - w.cachedInsets.Width
		dy := bounds.Dy() - w.cachedInsets.Height

		bounds.Min.X = w.insets.Left
		bounds.Min.Y = w.insets.Top
		bounds.Max.X = dx - w.insets.Right
		bounds.Max.Y = dy - w.insets.Bottom

		// Update bounds for the child
		w.cachedBounds = bounds
		w.child.SetBounds(bounds)
	}
}
func (w *tabsElement) updateProps(data *Tabs) error {
	return nil
}
