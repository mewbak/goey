package goey

import (
	"testing"
)

func (w *mountedVBox) Props() Widget {
	children := []Widget(nil)
	if len(w.children) != 0 {
		children = make([]Widget, 0, len(w.children))
		for _, v := range w.children {
			children = append(children, v.(Proper).Props())
		}
	}

	return &VBox{
		AlignMain:  w.alignMain,
		AlignCross: w.alignCross,
		Children:   children,
	}
}

func TestVBox(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingRenderWidgets(t,
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
		&VBox{Children: buttons, AlignMain: MainCenter},
		&VBox{Children: buttons, AlignMain: MainEnd},
		&VBox{Children: buttons, AlignMain: SpaceAround},
		&VBox{Children: buttons, AlignMain: SpaceBetween},
		&VBox{Children: buttons, AlignMain: Homogeneous},
	)
}

func TestVBoxClose(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t,
		&VBox{},
		&VBox{Children: buttons, AlignMain: MainStart},
	)
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
