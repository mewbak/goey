package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
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
	var values [3]bool

	testingCheckClick(t,
		&Checkbox{Text: "A", OnChange: func(v bool) { values[0] = v }},
		&Checkbox{Text: "B", Value: true, OnChange: func(v bool) { values[1] = v }},
		&Checkbox{Text: "C", OnChange: func(v bool) { values[2] = v }},
	)

	if !values[0] || values[1] || !values[2] {
		t.Errorf("OnChange failed, expected %v, got %v", [3]bool{true, false, true}, values[:])
	}
}

func TestCheckboxUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	}, []base.Widget{
		&Checkbox{Value: true, Text: "A--", Disabled: true},
		&Checkbox{Value: false, Text: "B--", Disabled: false},
	})
}
