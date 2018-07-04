package goey

import (
	"testing"
)

func TestEmptyCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Empty{},
		&Empty{},
		&Empty{},
	)
}

func TestEmptyClose(t *testing.T) {
	testingCloseWidgets(t,
		&Empty{},
		&Empty{},
	)
}
