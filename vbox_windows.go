package goey

import (
	"image"

	"github.com/lxn/win"
)

type mountedVBox struct {
	parent   NativeWidget
	children []MountedWidget
}

func (w *VBox) Mount(parent NativeWidget) (MountedWidget, error) {
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
}

func (w *mountedVBox) MeasureWidth() (DP, DP) {
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

func calculateGap(previous MountedWidget, current MountedWidget) DP {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// this is reduced to 5 relative pixels between a label and the following
	// control.  This relationship is not capture in the widget tree, so we
	// need to infer the relationship.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if _, ok := previous.(*mountedLabel); ok {
		return 5
	}
	if _, ok := previous.(*mountedCheckbox); ok {
		if _, ok := current.(*mountedCheckbox); ok {
			return 7
		}
	}

	return 11
}

func (w *mountedVBox) MeasureHeight(width DP) (DP, DP) {
	if len(w.children) == 0 {
		return 0, 0
	}

	previous := w.children[0]
	min, max := previous.MeasureHeight(width)
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureHeight(width)
		gap := calculateGap(previous, v)
		min += tmpMin + gap
		max += tmpMax + gap
		previous = v
	}
	return min, max
}

func (w *mountedVBox) SetBounds(bounds image.Rectangle) {
	if len(w.children) == 0 {
		return
	}

	posY := bounds.Min.Y
	width := bounds.Dx()
	widthDP := DP(width * 96 / dpi.X)

	previous := w.children[0]
	min, _ := previous.MeasureHeight(widthDP)
	height := min.ToPixelsY()
	previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
	posY += height

	// Assuming that height of bounds is sufficient
	for _, v := range w.children[1:] {
		posY += calculateGap(previous, v).ToPixelsY()
		previous = v

		min, _ := previous.MeasureHeight(widthDP)
		height := min.ToPixelsY()
		v.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
		posY += height
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
