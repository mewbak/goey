package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestTextArea(t *testing.T) {
	// Note, cannot use zero value for MinLines.  This will be changed to a
	// default value, and cause the post mounting check that the widget was
	// correctly instantiated to fail.
	testingRenderWidgets(t,
		&TextArea{Value: "A", MinLines: 3},
		&TextArea{Value: "B", MinLines: 3, Placeholder: "..."},
		&TextArea{Value: "C", MinLines: 3, Disabled: true},
	)
}

func TestTextAreaEvents(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&TextArea{},
		&TextArea{},
		&TextArea{},
	)
}

func TestTextAreaProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&TextArea{Value: "A", MinLines: 5},
		&TextArea{Value: "B", MinLines: 3, Placeholder: "..."},
		&TextArea{Value: "C", MinLines: 3, Disabled: true},
	}, []base.Widget{
		&TextArea{Value: "AA", MinLines: 6},
		&TextArea{Value: "BA", MinLines: 3, Disabled: true},
		&TextArea{Value: "CA", MinLines: 3, Placeholder: "***", Disabled: false},
	})
}
