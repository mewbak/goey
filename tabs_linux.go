// +build !gnustep

package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/gtk"
)

type tabsElement struct {
	handle   *gtk.Notebook
	value    int
	child    base.Element
	widgets  []TabItem
	insets   Insets
	onChange func(int)

	cachedInsets base.Point
	cachedBounds base.Rectangle
	cachedTabsW  base.Length
}

func tabsAppendChildren(handle *gtk.Notebook, children []TabItem) error {
	for _, v := range children {
		// Every tab needs some contents.  We will use a layout so that we can
		// custom layout of the controls.
		contents, err := gtk.LayoutNew(nil, nil)
		if err != nil {
			return err
		}

		// Create a label for the tabs.
		label, err := gtk.LabelNew(v.Caption)
		if err != nil {
			contents.Destroy()
			return err
		}

		// Append the new page to the notebook.
		handle.AppendPage(contents, label)
		contents.Show()
	}

	return nil
}

func (w *Tabs) mount(parent base.Control) (base.Element, error) {
	// Create the notebook control, and add tabs.
	control, err := gtk.NotebookNew()
	if err != nil {
		return nil, err
	}
	err = tabsAppendChildren(control, w.Children)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	parent.Handle.Add(control)

	child := base.Element(nil)
	if len(w.Children) > 0 {
		parent := getTabParent(control, w.Value)
		child_, err := base.Mount(parent, w.Children[w.Value].Child)
		if err != nil {
			control.Destroy()
			return nil, err
		}
		child = child_
		control.SetCurrentPage(w.Value)
	}

	retval := &tabsElement{
		handle:   control,
		child:    child,
		value:    w.Value,
		widgets:  w.Children,
		insets:   w.Insets,
		onChange: w.OnChange,
	}

	control.Connect("destroy", tabsOnDestroy, retval)
	control.Connect("switch-page", tabsOnSwitchPage, retval)
	control.ShowAll()

	return retval, nil
}

func tabsOnDestroy(widget *gtk.Notebook, mounted *tabsElement) {
	mounted.handle = nil
}

func tabsOnSwitchPage(widget *gtk.Notebook, _ *gtk.Widget, page int, mounted *tabsElement) {
	if page != mounted.value {
		if mounted.onChange != nil {
			mounted.onChange(page)
		}
		if page != mounted.value {
			mounted.mountPage(page)
			mounted.value = page
		}
	}
}

func getTabLabel(book *gtk.Notebook, ndx int) *gtk.Label {
	widget, err := book.GetNthPage(ndx)
	if err != nil {
		panic(err)
	}
	label, err := book.GetTabLabel(widget)
	if err != nil {
		panic(err)
	}
	return (*gtk.Label)(unsafe.Pointer(label))
}

func getTabParent(book *gtk.Notebook, ndx int) base.Control {
	widget, err := book.GetNthPage(ndx)
	if err != nil {
		panic(err)
	}
	layout := (*gtk.Layout)(unsafe.Pointer(widget))
	return base.Control{&layout.Container}
}

func (w *tabsElement) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *tabsElement) controlInsets() base.Point {
	if w.cachedInsets.Y == 0 {
		h1, _ := w.handle.GetPreferredHeight()
		// How should the offset between the notebook widget and the contained
		// page be measured?
		w.cachedInsets = base.Point{
			X: 0,
			Y: base.FromPixelsY(h1),
		}
	}

	return w.cachedInsets
}

func (w *tabsElement) controlTabsMinWidth() base.Length {
	if w.cachedTabsW == 0 {
		w1, _ := w.handle.GetPreferredWidth()
		w.cachedTabsW = base.FromPixelsX(w1)
	}
	return w.cachedTabsW
}

func (w *tabsElement) mountPage(page int) error {
	parent := getTabParent(w.handle, page)
	child, err := w.widgets[page].Child.Mount(parent)
	if err != nil {
		return err
	}
	child.Layout(base.Tight(base.Size{
		Width:  w.cachedBounds.Dx(),
		Height: w.cachedBounds.Dy(),
	}))
	child.SetBounds(w.cachedBounds)

	if w.child != nil {
		w.child.Close()
	}
	w.child = child
	return nil
}

func (w *tabsElement) Props() base.Widget {
	count := w.handle.GetNPages()
	children := make([]TabItem, count)
	for i := 0; i < count; i++ {
		label := getTabLabel(w.handle, i)
		text, _ := label.GetText()
		children[i].Caption = text
		children[i].Child = w.widgets[i].Child
	}

	return &Tabs{
		Value:    w.value,
		Children: children,
		OnChange: w.onChange,
	}
}

func (w *tabsElement) SetBounds(bounds base.Rectangle) {
	handle := Control{&w.handle.Widget}
	handle.SetBounds(bounds)

	if w.child != nil {
		// Determine the bounds for the child widget
		insets := w.controlInsets()
		insets.X += w.insets.Left + w.insets.Right
		insets.Y += w.insets.Top + w.insets.Bottom
		bounds = base.Rectangle{
			Max: base.Point{bounds.Dx() - insets.X, bounds.Dy() - insets.Y},
		}

		// Offset
		offset := base.Point{w.insets.Left, w.insets.Top}
		bounds.Min = bounds.Min.Add(offset)
		bounds.Max = bounds.Max.Add(offset)

		// Update bounds for the child
		w.cachedBounds = bounds
		w.child.SetBounds(bounds)
	}
}

func (w *tabsElement) updateProps(data *Tabs) error {
	if len(w.widgets) > len(data.Children) {
		// Modify captions for existing tabs
		for i, v := range data.Children {
			label := getTabLabel(w.handle, i)
			label.SetText(v.Caption)
		}
		// Remove excess tabs
		for i := len(w.widgets); i > len(data.Children); i-- {
			w.handle.RemovePage(i - 1)
		}
	} else {
		// Modify captions for existing tabs
		for i, v := range data.Children[:len(w.widgets)] {
			label := getTabLabel(w.handle, i)
			label.SetText(v.Caption)
		}
		// Append new tabs
		err := tabsAppendChildren(w.handle, data.Children[len(w.widgets):])
		if err != nil {
			return err
		}
	}
	w.widgets = data.Children

	// Update the selected widget
	if data.Value == w.value {
		w.mountPage(data.Value)
	} else {
		parent := getTabParent(w.handle, data.Value)
		child, err := base.DiffChild(parent, w.child, data.Children[data.Value].Child)
		w.child = child
		if err != nil {
			return err
		}

		w.handle.SetCurrentPage(data.Value)
		w.value = data.Value
	}

	return nil
}
