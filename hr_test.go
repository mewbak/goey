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
