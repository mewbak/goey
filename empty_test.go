package goey

import (
	"testing"
)

func TestEmpty(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Empty{},
		&Empty{},
		&Empty{},
	})
}
