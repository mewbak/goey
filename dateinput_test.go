package goey

import (
	"testing"
	"time"
)

func TestDateInputCreate(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingRenderWidgets(t, []Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	})
}

func TestDateInputClose(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingCloseWidgets(t, []Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	})
}

func TestDateInputEvents(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingCheckFocusAndBlur(t, []Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2},
		&DateInput{Value: v2},
	})
}

func TestDateInputUpdateProps(t *testing.T) {
	v1 := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	v2 := time.Date(2007, time.January, 2, 15, 4, 5, 0, time.Local)

	testingUpdateWidgets(t, []Widget{
		&DateInput{Value: v1},
		&DateInput{Value: v2, Disabled: true},
		&DateInput{Value: v2},
	}, []Widget{
		&DateInput{Value: v2},
		&DateInput{Value: v2, Disabled: false},
		&DateInput{Value: v1, Disabled: true},
	})
}
