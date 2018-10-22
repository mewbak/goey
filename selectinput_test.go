// +build !gnustep

package goey

import (
	"testing"
)

func TestSelectInputCreate(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingRenderWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
	)
}

func TestSelectInputClose(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingCloseWidgets(t,
		&SelectInput{Value: 0, Items: options},
		&SelectInput{Value: 1, Items: options},
		&SelectInput{Value: 2, Items: options, Disabled: true},
	)
}

func TestSelectInputEvents(t *testing.T) {
	options := []string{"Option A", "Option B", "Option C"}

	testingCheckFocusAndBlur(t,
		&SelectInput{Items: options},
		&SelectInput{Items: options},
		&SelectInput{Items: options},
	)
}
