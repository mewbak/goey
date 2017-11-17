package goey

import (
	"image"

	"github.com/lxn/win"
)

type MountedVBox struct {
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

	return &MountedVBox{parent: parent, children: c}, nil
}

func (w *MountedVBox) MinimumWidth() DP {
	if len(w.children) == 0 {
		return 0
	}

	retval := w.children[0].MinimumWidth()
	for _, v := range w.children[1:] {
		tmp := v.MinimumWidth()
		if tmp > retval {
			retval = tmp
		}
	}
	return retval
}

func calculateGap(previous MountedWidget, current MountedWidget) DP {
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

func (w *MountedVBox) CalculateHeight(width DP) DP {
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
	widthDP := DP(width * 96 / dpi.X)

	previous := w.children[0]
	height := previous.CalculateHeight(widthDP).ToPixelsY()
	previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
	posY += height

	// Assuming that height of bounds is sufficient
	for _, v := range w.children[1:] {
		posY += calculateGap(previous, v).ToPixelsY()
		previous = v

		height := v.CalculateHeight(widthDP).ToPixelsY()
		v.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+height))
		posY += height
	}
}

func (w *MountedVBox) SetChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.parent, w.children, children)
	return err
}

func (w *MountedVBox) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
