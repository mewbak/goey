package goey

import (
	"testing"
)

func TestLabel(t *testing.T) {
	testingRenderWidgets(t,
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	)
}

func TestLabelClose(t *testing.T) {
	testingCloseWidgets(t,
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	)
}

func TestLabelUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	}, []Widget{
		&Label{Text: "AB"},
		&Label{Text: "BC"},
		&Label{Text: "CD"},
	})
}
