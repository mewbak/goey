package goey

import (
	"image/color"
	"testing"
)

var (
	black = color.RGBA{0, 0, 0, 0xff}
	white = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func decorationChildWidget(child Element) Widget {
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
	testingUpdateWidgets(t, []Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
		&Decoration{Stroke: black},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
	}, []Widget{
		&Decoration{},
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{Fill: black, Stroke: white, Radius: 4 * DIP},
		&Decoration{Stroke: black},
	})
}
