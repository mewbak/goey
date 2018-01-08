package goey

import (
	"testing"
)

func TestParagraph(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&P{Text: "A", Align: Left},
		&P{Text: "B", Align: Right},
		&P{Text: "C", Align: Center},
		&P{Text: "D", Align: Justify},
	})
}
