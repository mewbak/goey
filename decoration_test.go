package goey

import (
	"image/color"
	"testing"

	"bitbucket.org/rj/goey/base"
)

var (
	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func (w *decorationElement) Props() base.Widget {
	widget := w.props()
	if w.child != nil {
		widget.Child = w.child.(Proper).Props()
	}

	return widget
}

func decorationChildWidget(child base.Element) base.Widget {
	if child == nil {
		return nil
	}

	return child.(Proper).Props()
}

func TestDecorationCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	)
}

func TestDecorationClose(t *testing.T) {
	testingCloseWidgets(t,
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	)
}

func TestDecorationUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	}, []base.Widget{
		&Decoration{},
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Stroke: black},
	})
}
