package goey

import (
	"github.com/lxn/win"
)

type mountedVBox struct {
	parent     NativeWidget
	children   []MountedWidget
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign
}

func (w *VBox) mount(parent NativeWidget) (MountedWidget, error) {
	c := make([]MountedWidget, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedVBox{parent: parent, children: c}, nil
}

func (w *mountedVBox) Close() {
	// On this platform, the mountedVBox handles layout, but does not actually
	// have an HWND, so there are no resources to release.

	// However, still need to free the children
	for _, v := range w.children {
		v.Close()
	}
	w.children = nil
}

func (w *mountedVBox) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	min, max := w.children[0].MeasureWidth()
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureWidth()
		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedVBox) MeasureHeight(width Length) (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	previous := w.children[0]
	min, max := previous.MeasureHeight(width)
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureHeight(width)
		gap := calculateVGap(previous, v)
		min += tmpMin + gap
		max += tmpMax + gap
		previous = v
	}
	return min, max
}

func (w *mountedVBox) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	width := bounds.Dx()
	minTotal, maxTotal := w.MeasureHeight(width)

	extraGap, deltaY, scale1, scale2 := distributeVSpace(w.alignMain, len(w.children), bounds.Dy(), minTotal, maxTotal)
	bounds.Min.Y += deltaY

	// Assuming that height of bounds is sufficient
	previous := MountedWidget(nil)
	for _, v := range w.children {
		if previous != nil {
			bounds.Min.Y += calculateVGap(previous, v) + extraGap
		}

		deltaY := setBoundsWithAlign(v, bounds, w.alignCross, scale1, scale2)
		bounds.Min.Y += deltaY
		previous = v
	}
}

func (w *mountedVBox) SetChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.parent, w.children, children)
	return err
}

func (w *mountedVBox) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}

func (w *mountedVBox) updateProps(data *VBox) error {
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	return w.SetChildren(data.Children)
}
