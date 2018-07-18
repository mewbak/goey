package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func (w *paddingElement) Props() base.Widget {
	child := base.Widget(nil)
	if w.child != nil {
		child = w.child.(Proper).Props()
	}

	return &Padding{
		Insets: w.insets,
		Child:  child,
	}
}

func TestPaddingCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)
}

func TestPaddingClose(t *testing.T) {
	testingCloseWidgets(t,
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	)
}

func TestPaddingUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Padding{Child: &Button{Text: "A"}},
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "B"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "C"}},
		&Padding{},
	}, []base.Widget{
		&Padding{Insets: DefaultInsets(), Child: &Button{Text: "AB"}},
		&Padding{Insets: UniformInsets(48 * DIP), Child: &Button{Text: "BC"}},
		&Padding{},
		&Padding{Child: &Button{Text: "CD"}},
	})
}
