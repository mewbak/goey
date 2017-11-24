package goey

import (
	"image"

	"github.com/lxn/win"
)

type mountedVBox struct {
	parent     NativeWidget
	children   []MountedWidget
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign
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

func (w *mountedVBox) MeasureWidth() (DIP, DIP) {
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

func calculateGap(previous MountedWidget, current MountedWidget) DIP {
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

func (w *mountedVBox) MeasureHeight(width DIP) (DIP, DIP) {
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
	widthDP := ToDIPX(width)
	height := bounds.Dy()
	heightDP := ToDIPY(height)
	minTotal, maxTotal := w.MeasureHeight(widthDP)

	if heightDP < minTotal {
		panic("not implemented")
	}

	// If there is more space than necessary, then we need to distribute the extra space.
	extraGap := 0
	if heightDP >= maxTotal {
		switch w.alignMain {
		case MainStart:
			// No need to do any adjustment.  The algorithm below will lay out
			// controls aligned to the top.
		case MainCenter:
			// Adjust the starting position to align the contents.
			posY += (height - maxTotal.PixelsY()) / 2

		case MainEnd:
			// Adjust the starting position to align the contents.
			posY += height - maxTotal.PixelsY()

		case SpaceAround:
			extraGap = (heightDP - maxTotal).PixelsY() / (len(w.children) + 1)
			posY += extraGap

		case SpaceBetween:
			if len(w.children) > 1 {
				extraGap = (heightDP - maxTotal).PixelsY() / (len(w.children) - 1)
			} else {
				// There are no controls between which to put the extra space.
				// The following essentially convert SpaceBetween to SpaceAround
				extraGap = (heightDP - maxTotal).PixelsY() / (len(w.children) + 1)
				posY += extraGap
			}
		}

		// Reduce available height
		heightDP = maxTotal
		height = heightDP.PixelsY()
	}

	scale1, scale2 := DIP(0), DIP(1)
	if heightDP > minTotal && maxTotal > minTotal {
		scale1, scale2 = heightDP-minTotal, maxTotal-minTotal
	}

	previous := w.children[0]
	min, max := previous.MeasureHeight(widthDP)
	h := (min + (max-min)*scale1/scale2).PixelsY()
	switch w.alignCross {
	case CrossStart:
		_, maxX := previous.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Min.X+newWidth, posY+h))
		} else {
			previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+h))
		}
	case CrossCenter:
		_, maxX := previous.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			x1 := (bounds.Min.X + bounds.Max.X - newWidth) / 2
			x2 := (bounds.Min.X + bounds.Max.X + newWidth) / 2
			previous.SetBounds(image.Rect(x1, posY, x2, posY+h))
		} else {
			previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+h))
		}
	case CrossEnd:
		_, maxX := previous.MeasureWidth()
		if newWidth := maxX.PixelsX(); newWidth < width {
			previous.SetBounds(image.Rect(bounds.Max.X-newWidth, posY, bounds.Max.X, posY+h))
		} else {
			previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+h))
		}
	case Stretch:
		previous.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+h))
	}
	posY += h

	// Assuming that height of bounds is sufficient
	for _, v := range w.children[1:] {
		posY += calculateGap(previous, v).PixelsY() + extraGap
		previous = v

		min, max := previous.MeasureHeight(widthDP)
		h := (min + (max-min)*scale1/scale2).PixelsY()
		v.SetBounds(image.Rect(bounds.Min.X, posY, bounds.Max.X, posY+h))
		posY += h
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

func (w *mountedVBox) UpdateProps(data_ Widget) error {
	data := data_.(*VBox)

	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	return w.SetChildren(data.Children)
}
