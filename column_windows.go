package goey

import (
	"github.com/lxn/win"
)

type mountedColumn struct {
	parent   NativeWidget
	children []Element
	counts   []int

	transition Length
}

func (w *Column) mount(parent NativeWidget) (Element, error) {
	c := make([]Element, 0, len(w.Children))
	counts := make([]int, 0, len(w.Children))

	for _, v := range w.Children {
		for _, w := range v {
			mountedChild, err := w.Mount(parent)
			if err != nil {
				return nil, err
			}
			c = append(c, mountedChild)
		}
		counts = append(counts, len(v))
	}

	return &mountedColumn{
		parent:   parent,
		children: c,
		counts:   counts,
	}, nil
}

func (w *mountedColumn) Close() {
	// On this platform, the mountedColumn handles layout, but does not actually
	// have an HWND, so there are no direct resources to release.

	// However, still need to free the children
	for _, v := range w.children {
		v.Close()
	}
	w.children = nil
}

func (w *mountedColumn) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	gap := calculateHGap(nil, nil)
	min, max := Length(0), Length(0)

	ndx := 0
	for _, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v
		tmpMin, tmpMax := vbox.MeasureWidth()

		if tmpMin > min {
			min = tmpMin
		}
		max = max + tmpMax + gap
	}
	w.transition = min.Scale(len(w.counts), 1) + gap.Scale(len(w.counts)-1, 1)
	println("transition", min.String(), len(w.counts), gap.String(), w.transition.String())
	return min, max
}

func (w *mountedColumn) MeasureHeight(width Length) (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	if w.transition == 0 {
		w.MeasureWidth()
		if w.transition == 0 {
			return 0, 0
		}
	}

	// If now side enough, we will layout the items exactly like a VBox
	if width < w.transition {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children,
		}

		return vbox.MeasureHeight(width)
	}

	ndx := 0
	min, max := Length(0), Length(0)
	for _, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v
		tmpMin, tmpMax := vbox.MeasureHeight(width)

		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedColumn) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	if w.transition == 0 {
		panic("internal error")
	}

	// If now side enough, we will layout the items exactly like a VBox
	if bounds.Dx() < w.transition {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children,
		}

		vbox.SetBounds(bounds)
		return
	}

	ndx := 0
	count := len(w.counts)
	gap := calculateHGap(nil, nil)
	bounds.Max.X += gap
	for i, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v

		minx := bounds.Min.X + bounds.Dx().Scale(i, count)
		maxx := bounds.Min.X + bounds.Dx().Scale(i+1, count) - gap
		vbox.SetBounds(Rectangle{Point{minx, bounds.Min.Y}, Point{maxx, bounds.Max.Y}})
	}
}

func (w *mountedColumn) SetChildren(children [][]Widget) error {
	// Flatten list
	c := make([]Widget, 0, len(children))
	for _, v := range children {
		c = append(c, v...)
	}

	err := error(nil)
	w.children, err = diffChildren(w.parent, w.children, c)
	return err
}

func (w *mountedColumn) SetOrder(previous win.HWND) win.HWND {
	for _, v := range w.children {
		previous = v.SetOrder(previous)
	}
	return previous
}
