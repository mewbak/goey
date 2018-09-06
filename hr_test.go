package goey

import (
	"bitbucket.org/rj/goey/base"
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

func TestHRUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&HR{},
		&HR{},
	}, []base.Widget{
		&HR{},
		&HR{},
	})
}
