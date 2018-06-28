package goey

import (
	"testing"
)

func TestDecorationCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
	})
}

func TestDecorationClose(t *testing.T) {
	testingCloseWidgets(t, []Widget{
		&Decoration{Child: &Button{Text: "A"}},
		&Decoration{},
	})
}
