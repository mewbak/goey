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

func TestParagraphProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&P{Text: "A", Align: Left},
		&P{Text: "B", Align: Right},
		&P{Text: "C", Align: Center},
		&P{Text: "D", Align: Justify},
	}, []Widget{
		&P{Text: "AAA", Align: Right},
		&P{Text: "BAA", Align: Center},
		&P{Text: "CAA", Align: Justify},
		&P{Text: "DAA", Align: Left},
	})
}
