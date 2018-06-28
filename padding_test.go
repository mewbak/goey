package goey

import (
	"testing"
)

func (w *paddingElement) Props() Widget {
	child, _ := w.child.(Proper)
	return &Padding{
		Insets: w.insets,
		Child:  child.Props(),
	}
}

func TestPaddingCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
	})
}

func TestPaddingUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
	}, []Widget{
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "AB"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "BC"}},
		&Padding{Child: &Button{Text: "CD"}},
	})
}
