package goey

import (
	"testing"
)

func TestCheckboxCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	})
}

func TestCheckboxEvents(t *testing.T) {
	testingCheckFocusAndBlur(t, []Widget{
		&Checkbox{Text: "A"},
		&Checkbox{Text: "B"},
		&Checkbox{Text: "C"},
	})
}
