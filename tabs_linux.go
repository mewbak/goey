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
	onChange func(int)

	cachedBounds base.Rectangle
	cachedSize   base.Size
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

func (w *tabsElement) mountPage(page int) error {
	parent := getTabParent(w.handle, page)
	child, err := w.widgets[page].Child.Mount(parent)
	if err != nil {
		return err
	}
	child.SetBounds(w.cachedBounds)

	if w.child != nil {
		w.child.Close()
	}
	w.child = child
	return nil
}

func (w *tabsElement) minSize() base.Size {
	if w.cachedSize.Height == 0 {
		w1, _ := w.handle.GetPreferredWidth()
		h1, _ := w.handle.GetPreferredHeight()
		// How should the offset between the notebook widget and the contained
		// page be measured?
		w.cachedSize = base.Size{
			Width:  base.FromPixelsX(w1),
			Height: base.FromPixelsY(h1),
		}
	}

	return w.cachedSize
}

func (w *tabsElement) Layout(bc base.Constraints) base.Size {
	tabsSize := w.minSize()

	if w.child == nil {
		return bc.Constrain(tabsSize)
	}

	size := w.child.Layout(bc.Inset(0, tabsSize.Height))
	return base.Size{
		Width:  max(size.Width, tabsSize.Width),
		Height: size.Height + tabsSize.Height,
	}
}

func (w *tabsElement) MinIntrinsicHeight(width base.Length) base.Length {
	tabsSize := w.minSize()

	if w.child == nil {
		return tabsSize.Height
	}

	if width == base.Inf {
		return w.child.MinIntrinsicHeight(base.Inf) + tabsSize.Height
	}

	return w.child.MinIntrinsicHeight(width)
}

func (w *tabsElement) MinIntrinsicWidth(height base.Length) base.Length {
	tabsSize := w.minSize()

	if w.child == nil {
		return tabsSize.Width
	}

	if height == base.Inf {
		return w.child.MinIntrinsicWidth(base.Inf)
	}

	return w.child.MinIntrinsicWidth(height - tabsSize.Height)
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
		tabsSize := w.minSize()
		w.cachedBounds = base.Rectangle{
			Max: base.Point{bounds.Dx(), bounds.Dy() - tabsSize.Height},
		}

		// Update bounds for the child
		w.child.SetBounds(w.cachedBounds)
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
