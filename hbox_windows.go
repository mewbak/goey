package goey

import (
	"github.com/lxn/win"
	"image"
)

func (w *MountedHBox) PreferredWidth() int {
	if len(w.children) == 0 {
		return 0
	}

	retval := w.children[0].PreferredWidth()
	for _, v := range w.children[1:] {
		retval = retval + v.PreferredWidth() + 8
	}
	w.preferredWidth = retval
	return retval
}

func (w *MountedHBox) CalculateHeight(width int) int {
	if len(w.children) == 0 {
		return 0
	}

	if w.preferredWidth == 0 {
		w.PreferredWidth()
		if w.preferredWidth == 0 {
			return 0
		}
	}

	if w.preferredWidth >= width || w.align == Justify {
		width = (width + 8) / len(w.children)

		retval := w.children[0].CalculateHeight(width)
		for _, v := range w.children[1:] {
			tmp := v.CalculateHeight(width)
			if tmp > retval {
				retval = tmp
			}
		}
		return retval
	}

	retval := w.children[0].CalculateHeight(w.children[0].PreferredWidth())
	for _, v := range w.children[1:] {
		tmp := v.CalculateHeight(w.children[0].PreferredWidth())
		if tmp > retval {
			retval = tmp
		}
	}
	return retval
}

func (w *MountedHBox) SetBounds(bounds image.Rectangle) {
	width := bounds.Dx()
	length := len(w.children)

	if w.preferredWidth == 0 {
		w.PreferredWidth()
		if w.preferredWidth == 0 {
			return
		}
	}

	// Assuming that height of bounds is sufficient
	if w.preferredWidth >= width || w.align == Justify {
		for i, v := range w.children {
			posX1 := bounds.Min.X + (width+8)*i/length
			posX2 := bounds.Min.X + (width+8)*(i+1)/length - 8
			v.SetBounds(image.Rect(posX1, bounds.Min.Y, posX2, bounds.Max.Y))
		}
	} else if w.align == Left {
		posX := bounds.Min.X
		for _, v := range w.children {
			posX2 := posX + v.PreferredWidth()
			v.SetBounds(image.Rect(posX, bounds.Min.Y, posX2, bounds.Max.Y))
			posX = posX2 + 8
		}
	} else {
		posX := bounds.Max.X
		for i := len(w.children); i > 0; i-- {
			v := w.children[i-1]
			posX2 := posX - v.PreferredWidth()
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
