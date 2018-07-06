package goey

import (
	"testing"
)

func (w *hboxElement) Props() Widget {
	children := []Widget(nil)
	if len(w.children) != 0 {
		children = make([]Widget, 0, len(w.children))
		for _, v := range w.children {
			children = append(children, v.(Proper).Props())
		}
	}

	return &HBox{
		AlignMain:  w.alignMain,
		AlignCross: w.alignCross,
		Children:   children,
	}
}

func TestHBox(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingRenderWidgets(t,
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
		&HBox{Children: buttons, AlignMain: MainCenter},
		&HBox{Children: buttons, AlignMain: MainEnd},
		&HBox{Children: buttons, AlignMain: SpaceAround},
		&HBox{Children: buttons, AlignMain: SpaceBetween},
		&HBox{Children: buttons, AlignMain: Homogeneous},
	)
}

func TestHBoxClose(t *testing.T) {
	buttons := []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	}

	testingCloseWidgets(t,
		&HBox{},
		&HBox{Children: buttons, AlignMain: MainStart},
	)
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

func TestHBoxLayout(t *testing.T) {
	cases := []struct {
		children           []Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicWidth  Length
		minIntrinsicHeight Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{[]Element{&mockElement{13 * DIP, 13 * DIP}, &mockElement{13 * DIP, 13 * DIP}}, MainStart, Stretch, 37 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 13 * DIP}, &mockElement{13 * DIP, 15 * DIP}}, MainStart, Stretch, 37 * DIP, 15 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, MainStart, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{26 * DIP, 13 * DIP}, &mockElement{13 * DIP, 11 * DIP}}, Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
	}

	for i, v := range cases {
		in := hboxElement{
			children:   v.children,
			alignMain:  v.alignMain,
			alignCross: v.alignCross,
		}

		if value := in.MinIntrinsicHeight(Inf); value != v.minIntrinsicHeight {
			t.Errorf("Incorrect min intrinsic height on case %d, got %s, want %s", i, value, v.minIntrinsicHeight)
		}
		if value := in.MinIntrinsicWidth(Inf); value != v.minIntrinsicWidth {
			t.Errorf("Incorrect min intrinsic width on case %d, got %s, want %s", i, value, v.minIntrinsicWidth)
		}
	}
}
