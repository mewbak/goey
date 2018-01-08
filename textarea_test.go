package goey

import (
	"testing"
)

func TestTextArea(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&TextArea{Value: "A"},
		&TextArea{Value: "B", Placeholder: "..."},
		&TextArea{Value: "C", Disabled: true},
	})
}

func TestTextAreaEvents(t *testing.T) {
	testingCheckFocusAndBlur(t, []Widget{
		&TextArea{},
		&TextArea{},
		&TextArea{},
	})
}
