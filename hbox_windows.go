package goey

import (
	"image"

	"github.com/lxn/win"
)

type mountedHBox struct {
	parent   NativeWidget
	children []MountedWidget
	align    TextAlignment

	minimumWidth DIP
}

func (w *HBox) Mount(parent NativeWidget) (MountedWidget, error) {
	c := make([]MountedWidget, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	align := w.Align
	if align == DefaultAlign {
		align = Justify
	}

	return &mountedHBox{parent: parent, children: c, align: align}, nil
}

func (w *mountedHBox) Close() {
	// nothing required
}

func (w *mountedHBox) MeasureWidth() (DIP, DIP) {
	if len(w.children) == 0 {
		return 0, 0
	}

	min, max := w.children[0].MeasureWidth()
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureWidth()
		min = min + tmpMin + 8
		max = max + tmpMax + 8
	}
	w.minimumWidth = min
	return min, max
}

func (w *mountedHBox) MeasureHeight(width DIP) (DIP, DIP) {
	if len(w.children) == 0 {
		return 0, 0
	}

	if w.minimumWidth == 0 {
		w.MeasureWidth()
		if w.minimumWidth == 0 {
			return 0, 0
		}
	}

	if w.minimumWidth >= width || w.align == Justify {
		width = (width + 8) / DIP(len(w.children))

		min, max := w.children[0].MeasureHeight(width)
		for _, v := range w.children[1:] {
			tmpMin, tmpMax := v.MeasureHeight(width)
			if tmpMin > min {
				min = tmpMin
			}
			if tmpMax > max {
				max = tmpMax
			}
		}
		return min, max
	}

	minWidth, _ := w.children[0].MeasureWidth()
	min, max := w.children[0].MeasureHeight(minWidth)
	for _, v := range w.children[1:] {
		minWidth, _ = v.MeasureWidth()
		tmpMin, tmpMax := v.MeasureHeight(minWidth)
		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedHBox) SetBounds(bounds image.Rectangle) {
	width := bounds.Dx()
	widthDP := ToDIPX(width)
	length := len(w.children)

	if w.minimumWidth == 0 {
		w.MeasureWidth()
		if w.minimumWidth == 0 {
			return
		}
	}

	// Assuming that height of bounds is sufficient
	if w.minimumWidth >= widthDP || w.align == Justify {
		for i, v := range w.children {
			posX1 := bounds.Min.X + (width+8)*i/length
			posX2 := bounds.Min.X + (width+8)*(i+1)/length - 8
			v.SetBounds(image.Rect(posX1, bounds.Min.Y, posX2, bounds.Max.Y))
		}
	} else if w.align == Left {
		posX := bounds.Min.X
		for _, v := range w.children {
			min, _ := v.MeasureWidth()
			posX2 := posX + min.PixelsX()
			v.SetBounds(image.Rect(posX, bounds.Min.Y, posX2, bounds.Max.Y))
			posX = posX2 + 8
		}
	} else {
		posX := bounds.Max.X
		for i := len(w.children); i > 0; i-- {
			v := w.children[i-1]
			min, _ := v.MeasureWidth()
			posX2 := posX - min.PixelsX()
			v.SetBounds(image.Rect(posX2, bounds.Min.Y, posX, bounds.Max.Y))
			posX = posX2 - 8
		}
	}
}

func (w *mountedHBox) SetChildren(children []Widget) error {
	panic("not implemented")
}

func (w *mountedHBox) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
