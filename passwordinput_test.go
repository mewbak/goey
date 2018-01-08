package goey

import (
	"testing"
)

func TestPasswordInput(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&PasswordInput{Value: "A"},
		&PasswordInput{Value: "B", Placeholder: "..."},
		&PasswordInput{Value: "C", Disabled: true},
	})
}

func TestPasswordInputEvents(t *testing.T) {
	testingCheckFocusAndBlur(t, []Widget{
		&PasswordInput{},
		&PasswordInput{},
		&PasswordInput{},
	})
}
