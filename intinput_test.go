package goey

import (
	"testing"
)

func TestIntInput(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	})
}

func TestIntInputEvents(t *testing.T) {
	testingCheckFocusAndBlur(t, []Widget{
		&IntInput{},
		&IntInput{},
		&IntInput{},
	})
}
