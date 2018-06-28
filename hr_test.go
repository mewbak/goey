package goey

import (
	"testing"
)

func TestHR(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&HR{},
		&HR{},
		&HR{},
	})
}

func TestHRClose(t *testing.T) {
	testingCloseWidgets(t, []Widget{
		&HR{},
		&HR{},
	})
}
