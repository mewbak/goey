package goey

import (
	"testing"
)

func TestHBox(t *testing.T) {
buttons := []Widget{
	&Button{Text:"A"},
	&Button{Text:"B"},
	&Button{Text:"C"},
}

	testingRenderWidgets(t, []Widget{
		&HBox{},
		&HBox{Children:buttons,AlignMain:MainStart},
		&HBox{Children:buttons,AlignMain:MainCenter},
		&HBox{Children:buttons,AlignMain:MainEnd},
		&HBox{Children:buttons,AlignMain:SpaceAround},
		&HBox{Children:buttons,AlignMain:SpaceBetween},
	})
}
