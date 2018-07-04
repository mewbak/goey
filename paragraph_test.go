package goey

import (
	"testing"
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
	testingUpdateWidgets(t, []Widget{
		&P{Text: "A", Align: JustifyLeft},
		&P{Text: "B", Align: JustifyRight},
		&P{Text: "C", Align: JustifyCenter},
		&P{Text: "D", Align: JustifyFull},
	}, []Widget{
		&P{Text: "AAA", Align: JustifyRight},
		&P{Text: "BAA", Align: JustifyCenter},
		&P{Text: "CAA", Align: JustifyFull},
		&P{Text: "DAA", Align: JustifyLeft},
	})
}
