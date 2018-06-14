package main

import (
	"bitbucket.org/rj/goey"
)

var (
	gridKind = goey.NewKind("bitbucket.org/rj/example/feettometer.Grid")
)

type Grid struct {
	Children [3][3]goey.Widget
	Gutter   goey.Length
}

func (*Grid) Kind() *goey.Kind {
	return &gridKind
}

func (w *Grid) Mount(parent goey.Control) (goey.Element, error) {
	retval := &mountedGrid{parent: parent, gutter: w.Gutter}
	if retval.gutter == 0 {
		retval.gutter = 11 * goey.DIP
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			w, err := w.Children[i][j].Mount(parent)
			if err != nil {
				retval.Close()
				return nil, err
			}
			retval.children[i][j] = w
		}
	}

	return retval, nil
}

type mountedGrid struct {
	parent   goey.Control
	gutter   goey.Length
	children [3][3]goey.Element
	heights  struct {
		min [3]goey.Length
		max [3]goey.Length
	}
}

func (w *mountedGrid) Close() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if child := w.children[i][j]; child != nil {
				child.Close()
				w.children[i][j] = nil
			}
		}
	}
}

func (w *mountedGrid) Kind() *goey.Kind {
	return &gridKind
}

func (w *mountedGrid) MeasureWidth() (goey.Length, goey.Length) {
	min, max := goey.Length(0), goey.Length(0)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			tmpMin, tmpMax := w.children[i][j].MeasureWidth()
			if tmpMin > min {
				min = tmpMin
			}
			if tmpMax > max {
				max = tmpMax
			}
		}
	}

	return min*3 + 2*w.gutter, max*3 + 2*w.gutter
}

func (w *mountedGrid) MeasureHeight(width goey.Length) (goey.Length, goey.Length) {
	min, max := goey.Length(0), goey.Length(0)
	for i := 0; i < 3; i++ {
		rowMin, rowMax := w.children[i][0].MeasureHeight(width / 3)
		for j := 1; j < 3; j++ {
			tmpMin, tmpMax := w.children[i][j].MeasureHeight(width / 3)
			if tmpMin > rowMin {
				rowMin = tmpMin
			}
			if tmpMax > rowMax {
				rowMax = tmpMax
			}
		}

		w.heights.min[i] = rowMin
		w.heights.max[i] = rowMax
		min += rowMin
		max += rowMax
	}

	return min + 2*w.gutter, max + 2*w.gutter
}

func (w *mountedGrid) SetBounds(bounds goey.Rectangle) {
	// We shoudl be smarter, but for the moment we just divide the rows and columns in three
	dx := bounds.Dx() + w.gutter
	dy := bounds.Dy() + w.gutter

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			x1 := bounds.Min.X + dx.Scale(j, 3)
			y1 := bounds.Min.Y + dy.Scale(i, 3)
			x2 := bounds.Min.X + dx.Scale(j+1, 3) - w.gutter
			y2 := bounds.Min.Y + dy.Scale(i+1, 3) - w.gutter
			w.children[i][j].SetBounds(goey.Rectangle{goey.Point{x1, y1}, goey.Point{x2, y2}})
		}
	}
}

func (w *mountedGrid) updateProps(data *Grid) error {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			child, err := goey.DiffChild(w.parent, w.children[i][j], data.Children[i][j])
			if err != nil {
				return err
			}
			w.children[i][j] = child
		}
	}
	return nil
}

func (w *mountedGrid) UpdateProps(data goey.Widget) error {
	return w.updateProps(data.(*Grid))
}
