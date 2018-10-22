// +build !gnustep

package goey

import (
	"testing"

	"bitbucket.org/rj/goey/base"
)

func TestProgressCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 0},
	)
}

func TestProgressClose(t *testing.T) {
	testingCloseWidgets(t,
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 0},
	)
}

func TestProgressUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Progress{Value: 50, Min: 0, Max: 100},
		&Progress{Value: 50, Min: 0, Max: 100},
	}, []base.Widget{
		&Progress{Value: 75, Min: 0, Max: 100},
		&Progress{Value: 50, Min: 0, Max: 200},
	})
}
