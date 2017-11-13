package goey

import (
	"github.com/lxn/win"
	"image"
)

func (w *MountedVBox) PreferredWidth() int {
	if len(w.children) == 0 {
		return 0
	}

	retval := w.children[0].PreferredWidth()
	for _, v := range w.children[1:] {
		tmp := v.PreferredWidth()
		if tmp > retval {
			retval = tmp
		}
	}
	return retval
}

func calculateGap(previous MountedWidget, current MountedWidget) int {
	// The vertical gap between most controls is 11 relative pixels.  However,
	// this is reduced to 5 relative pixels between a label and the following
	// control.  This relationship is not capture in the widget tree, so we
	// need to infer the relationship.
	//
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if _, ok := previous.(*MountedLabel); ok {
		return 5
	}
	if _, ok := previous.(*MountedCheckbox); ok {
		if _, ok := current.(*MountedCheckbox); ok {
			return 7
		}
	}

	return 11
}

func (w *MountedVBox) CalculateHeight(width int) int {
	if len(w.children) == 0 {
		return 0
	}

	previous := w.children[0]
	retval := previous.CalculateHeight(width)
	for _, v := range w.children[1:] {
		retval += calculateGap(previous, v) + v.CalculateHeight(width)
		previous = v
	}
	return retval
}

func (w *MountedVBox) SetBounds(bounds image.Rectangle) {
	if len(w.children) == 0 {
		return
	}

	posY := bounds.Min.Y
	width := bounds.Dx()

	previous := w.children[0]
	height := previous.CalculateHeight(width)
	previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
	posY += height

	// Assuming that height of bounds is sufficient
	for _, v := range w.children[1:] {
		posY += calculateGap(previous, v)
		previous = v

		height := v.CalculateHeight(width)
		v.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
		posY += height
	}
}

func (w *MountedVBox) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
