package goey

import (
	"testing"
)

func TestHBox(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingRenderWidgets(t, []Widget{
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
		&HBox{Children: buttons, AlignMain: MainCenter},
		&HBox{Children: buttons, AlignMain: MainEnd},
		&HBox{Children: buttons, AlignMain: SpaceAround},
		&HBox{Children: buttons, AlignMain: SpaceBetween},
		&HBox{Children: buttons, AlignMain: Homogeneous},
	})
}

func TestHBoxClose(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t, []Widget{
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
	})
}

func TestHBoxUpdateProps(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingUpdateWidgets(t, []Widget{
		&HBox{AlignMain: MainStart},
		&HBox{Children: buttons, AlignMain: MainEnd, AlignCross: CrossStart},
	}, []Widget{
		&HBox{Children: buttons, AlignMain: MainEnd},
		&HBox{AlignMain: MainStart, AlignCross: CrossCenter},
	})
}
