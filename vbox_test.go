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
	children := []Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}

	cases := []struct {
		children    []Element
		alignMain   MainAxisAlign
		alignCross  CrossAxisAlign
		constraints Constraint
		size        Size
		bounds      []Rectangle
	}{
		{nil, MainStart, Stretch, TightWidth(40 * DIP), Size{40 * DIP, 0}, []Rectangle{}},
		{children, MainStart, Stretch, TightWidth(40 * DIP), Size{40 * DIP, 50 * DIP}, []Rectangle{
			Rect(0, 0, 40*DIP, 26*DIP), Rect(0, 37*DIP, 40*DIP, 50*DIP),
		}},
	}

	for i, v := range cases {
		in := vboxElement{
			children:     v.children,
			alignMain:    v.alignMain,
			alignCross:   v.alignCross,
			childrenSize: make([]Size, len(v.children)),
		}

		size := in.Layout(v.constraints)
		if size != v.size {
			t.Errorf("Incorrect size on case %d, got %s, want %s", i, size, v.size)
		}
		in.SetBounds(Rectangle{Point{}, Point{size.Width, size.Height}})
		for j, u := range v.bounds {
			if got := v.children[j].(*mockElement).Bounds; got != u {
				t.Errorf("Incorrect bounds case %d-%d, got %s, want %s", i, j, got, u)
			}
		}
	}
}

func TestVBoxMinIntrinsic(t *testing.T) {
	cases := []struct {
		children           []Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicHeight Length
		minIntrinsicWidth  Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{[]Element{mock(13*DIP, 13*DIP), mock(13*DIP, 13*DIP)}, MainStart, Stretch, 37 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 13*DIP), mock(15*DIP, 13*DIP)}, MainStart, Stretch, 37 * DIP, 15 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, MainStart, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 26*DIP), mock(11*DIP, 13*DIP)}, Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
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
