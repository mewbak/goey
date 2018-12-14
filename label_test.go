package goey

import (
	"testing"
	"testing/quick"
	"unicode"

	"bitbucket.org/rj/goey/base"
)

func TestLabel(t *testing.T) {
	testingRenderWidgets(t,
		&Label{Text: "A"},
		&Label{Text: "B"},
		&Label{Text: "C"},
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
	)

	f := func(text string) bool {
		// Filter out bad unicode code points
		for _, v := range []rune(text) {
			if !unicode.IsGraphic(v) {
				return true
			}
		}
		return testingRenderWidget(t, &Label{Text: text})
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("quick: %s", err)
	}
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
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
	}, []base.Widget{
		&Label{Text: ""},
		&Label{Text: "ABCD\nEDFG"},
		&Label{Text: "AB"},
		&Label{Text: "BC"},
		&Label{Text: "CD"},
	})
}
