package goey

import (
	"testing"
	"testing/quick"
	"unicode"

	"bitbucket.org/rj/goey/base"
)

func TestParagraph(t *testing.T) {
	testingRenderWidgets(t,
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
	)

	f := func(text string, align uint) bool {
		// Filter out bad unicode code points
		for _, v := range []rune(text) {
			if !unicode.IsGraphic(v) {
				return true
			}
		}
		return testingRenderWidget(t, &P{Text: text, Align: TextAlignment(align % 4)})
	}
	if err := quick.Check(f, nil); err != nil {
		t.Errorf("quick: %s", err)
	}
}

func TestParagraphClose(t *testing.T) {
	testingCloseWidgets(t,
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
	)
}

func TestParagraphProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
	}, []base.Widget{
		&P{Text: "", Align: JustifyLeft},
		&P{Text: "ABCD\nEFGH", Align: JustifyLeft},
		&P{Text: "AAA", Align: JustifyRight},
		&P{Text: "BAA", Align: JustifyCenter},
		&P{Text: "CAA", Align: JustifyFull},
		&P{Text: "DAA", Align: JustifyLeft},
	})
}
