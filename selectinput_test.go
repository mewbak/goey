package goey

import (
	"testing"
)

func TestSelectInput(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingRenderWidgets(t, []Widget{
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
	})
}

func TestSelectInputEvents(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingCheckFocusAndBlur(t, []Widget{
		&SelectInput{Items: options},
		&SelectInput{Items: options},
		&SelectInput{Items: options},
	})
}
