package goey

import (
	"testing"
)

func TestCheckboxCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	)
}

func TestCheckboxClose(t *testing.T) {
	testingCloseWidgets(t,
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	)
}

func TestCheckboxFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&Checkbox{Text: "A"},
		&Checkbox{Text: "B"},
		&Checkbox{Text: "C"},
	)
}

func TestCheckboxClick(t *testing.T) {
	testingCheckClick(t,
		&Checkbox{Text: "A"},
		&Checkbox{Text: "B"},
		&Checkbox{Text: "C"},
	)
}

func TestCheckboxUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	}, []Widget{
		&Checkbox{Value: true, Text: "A--", Disabled: true},
		&Checkbox{Value: false, Text: "B--", Disabled: false},
	})
}
