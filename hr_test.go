package goey

import (
	"testing"
	"bitbucket.org/rj/goey/base"
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

func TestHRUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&HR{},
		&HR{},
	}, []base.Widget{
		&HR{},
		&HR{},
	})
}
