package goey

import (
	"testing"
)

func (w *alignElement) Props() Widget {
	child, _ := w.child.(Proper)
	return &Align{
		HAlign:       w.halign,
		VAlign:       w.valign,
		WidthFactor:  w.widthFactor,
		HeightFactor: w.heightFactor,
		Child:        child.Props(),
	}
}

func TestAlignCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
	})
}

func TestAlignUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&Align{Child: &Button{Text: "A"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "B"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "C"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "C"}},
		&Align{HeightFactor: 2, WidthFactor: 2.5, Child: &Button{Text: "C"}},
	}, []Widget{
		&Align{Child: &Button{Text: "AB"}},
		&Align{HAlign: AlignCenter, Child: &Button{Text: "BC"}},
		&Align{HAlign: AlignStart, Child: &Button{Text: "CD"}},
		&Align{HAlign: AlignEnd, Child: &Button{Text: "CE"}},
		&Align{HeightFactor: 4, WidthFactor: 3, Child: &Button{Text: "CF"}},
	})
}
