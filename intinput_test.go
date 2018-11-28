package goey

import (
	"bitbucket.org/rj/goey/base"
	"testing"
)

func TestIntInputMount(t *testing.T) {
	testingMountWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	)
}

func TestIntInputClose(t *testing.T) {
	testingCloseWidgets(t,
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	)
}

func TestIntInputFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&IntInput{},
		&IntInput{},
		&IntInput{},
	)
}

func TestIntInputUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 2, Placeholder: "..."},
		&IntInput{Value: 3, Disabled: true},
	}, []base.Widget{
		&IntInput{Value: 1},
		&IntInput{Value: 4, Disabled: true},
		&IntInput{Value: 5, Placeholder: "***"},
	})
}
