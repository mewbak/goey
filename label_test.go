package goey

import (
	"testing"
)

func TestLabel(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	})
}

func TestLabelClose(t *testing.T) {
	testingCloseWidgets(t, []Widget{
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	})
}
