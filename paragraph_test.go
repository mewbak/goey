// +build !gnustep

package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestParagraph(t *testing.T) {
	testingRenderWidgets(t,
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
	)
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
	}, []base.Widget{
		&P{Text: "AAA", Align: JustifyRight},
		&P{Text: "BAA", Align: JustifyCenter},
		&P{Text: "CAA", Align: JustifyFull},
		&P{Text: "DAA", Align: JustifyLeft},
	})
}
