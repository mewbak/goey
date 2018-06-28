package goey

import (
	"testing"
)

func TestVBox(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingRenderWidgets(t, []Widget{
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
		&VBox{Children: buttons, AlignMain: MainCenter},
		&VBox{Children: buttons, AlignMain: MainEnd},
		&VBox{Children: buttons, AlignMain: SpaceAround},
		&VBox{Children: buttons, AlignMain: SpaceBetween},
		&VBox{Children: buttons, AlignMain: Homogeneous},
	})
}

func TestVBoxClose(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t, []Widget{
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
	})
}

func TestVBoxUpdateProps(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingUpdateWidgets(t, []Widget{
		&VBox{AlignMain: MainStart},
		&VBox{Children: buttons, AlignMain: MainEnd, AlignCross: CrossStart},
	}, []Widget{
		&VBox{Children: buttons, AlignMain: MainEnd},
		&VBox{AlignMain: MainStart, AlignCross: CrossCenter},
	})
}
