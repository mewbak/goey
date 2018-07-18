package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func (w *alignElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Align{
		HAlign:       w.halign,
		VAlign:       w.valign,
		WidthFactor:  w.widthFactor,
		HeightFactor: w.heightFactor,
		Child:        child,
	}
}

func TestAlignCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
		&Align{},
	)
}

func TestAlignClose(t *testing.T) {
	testingCloseWidgets(t,
		&Align{Child: &Button{Text: "A"}},
		&Align{},
	)
}

func TestAlignUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
		&Align{},
	}, []base.Widget{
		&Align{Child: &Button{Text: "AB"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "BC"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "CD"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "CE"}},
		&Align{HeightFactor: 4, WidthFactor: 3},
		&Align{Child: &Label{Text: "CF"}},
	})
}
