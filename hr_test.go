package goey

import (
	"testing"
)

func TestHR(t *testing.T) {
	testingRenderWidgets(t,
		&HR{},
		&HR{},
		&HR{},
	)
}

func TestHRClose(t *testing.T) {
	testingCloseWidgets(t,
		&HR{},
		&HR{},
	)
}
