package goey

import (
	"image"

	"github.com/lxn/win"
)

type MountedHBox struct {
	parent   NativeWidget
	children []MountedWidget
	align    Alignment

	minimumWidth DP
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

	return &MountedHBox{parent: parent, children: c, align: align}, nil
}

func (w *MountedHBox) MinimumWidth() DP {
	if len(w.children) == 0 {
		return 0
	}

	retval := w.children[0].MinimumWidth()
	for _, v := range w.children[1:] {
		retval = retval + v.MinimumWidth() + 8
	}
	w.minimumWidth = retval
	return retval
}

func (w *MountedHBox) CalculateHeight(width DP) DP {
	if len(w.children) == 0 {
		return 0
	}

	if w.minimumWidth == 0 {
		w.MinimumWidth()
		if w.minimumWidth == 0 {
			return 0
		}
	}

	if w.minimumWidth >= width || w.align == Justify {
		width = (width + 8) / DP(len(w.children))

		retval := w.children[0].CalculateHeight(width)
		for _, v := range w.children[1:] {
			tmp := v.CalculateHeight(width)
			if tmp > retval {
				retval = tmp
			}
		}
		return retval
	}

	retval := w.children[0].CalculateHeight(w.children[0].MinimumWidth())
	for _, v := range w.children[1:] {
		tmp := v.CalculateHeight(w.children[0].MinimumWidth())
		if tmp > retval {
			retval = tmp
		}
	}
	return retval
}

func (w *MountedHBox) SetBounds(bounds image.Rectangle) {
	width := bounds.Dx()
	widthDP := DP(width * dpi.X / 96)
	length := len(w.children)

	if w.minimumWidth == 0 {
		w.MinimumWidth()
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
			posX2 := posX + v.MinimumWidth().ToPixelsX()
			v.SetBounds(image.Rect(posX, bounds.Min.Y, posX2, bounds.Max.Y))
			posX = posX2 + 8
		}
	} else {
		posX := bounds.Max.X
		for i := len(w.children); i > 0; i-- {
			v := w.children[i-1]
			posX2 := posX - v.MinimumWidth().ToPixelsX()
			v.SetBounds(image.Rect(posX2, bounds.Min.Y, posX, bounds.Max.Y))
			posX = posX2 - 8
		}
	}
}

func (w *MountedHBox) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
