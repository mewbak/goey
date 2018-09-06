package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
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
	testingUpdateWidgets(t, []base.Widget{
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
	}, []base.Widget{
		&Label{Text: "AB"},
		&Label{Text: "BC"},
		&Label{Text: "CD"},
	})
}
