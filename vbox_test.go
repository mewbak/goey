package goey

import (
	"testing"
)

func (w *vboxElement) Props() Widget {
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

func TestVBoxLayout(t *testing.T) {
	cases := []struct {
		children           []Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicHeight Length
		minIntrinsicWidth  Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{[]Element{&mockElement{13 * DIP, 13 * DIP}, &mockElement{13 * DIP, 13 * DIP}}, MainStart, Stretch, 37 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 13 * DIP}, &mockElement{15 * DIP, 13 * DIP}}, MainStart, Stretch, 37 * DIP, 15 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, MainStart, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{&mockElement{13 * DIP, 26 * DIP}, &mockElement{11 * DIP, 13 * DIP}}, Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
	}

	for i, v := range cases {
		in := vboxElement{
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
