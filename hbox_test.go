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
	children := []Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}

	cases := []struct {
		children    []Element
		alignMain   MainAxisAlign
		alignCross  CrossAxisAlign
		constraints Constraint
		size        Size
		bounds      []Rectangle
	}{
		{nil, MainStart, Stretch, TightHeight(40 * DIP), Size{0, 40 * DIP}, []Rectangle{}},
		{children, MainStart, Stretch, TightHeight(40 * DIP), Size{50 * DIP, 40 * DIP}, []Rectangle{
			Rect(0, 0, 26*DIP, 40*DIP), Rect(37*DIP, 0, 50*DIP, 40*DIP),
		}},
		{children, MainStart, Stretch, Tight(Size{150 * DIP, 40 * DIP}), Size{150 * DIP, 40 * DIP}, []Rectangle{
			Rect(0, 0, 26*DIP, 40*DIP), Rect(37*DIP, 0, 50*DIP, 40*DIP),
		}},
		{children, MainEnd, Stretch, Tight(Size{150 * DIP, 40 * DIP}), Size{150 * DIP, 40 * DIP}, []Rectangle{
			Rect(100*DIP, 0, 126*DIP, 40*DIP), Rect(137*DIP, 0, 150*DIP, 40*DIP),
		}},
	}

	for i, v := range cases {
		in := hboxElement{
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

func TestHBoxMinIntrinsic(t *testing.T) {
	cases := []struct {
		children           []Element
		alignMain          MainAxisAlign
		alignCross         CrossAxisAlign
		minIntrinsicWidth  Length
		minIntrinsicHeight Length
	}{
		{nil, MainStart, Stretch, 0, 0},
		{[]Element{mock(13*DIP, 13*DIP), mock(13*DIP, 13*DIP)}, MainStart, Stretch, 37 * DIP, 13 * DIP},
		{[]Element{mock(13*DIP, 13*DIP), mock(13*DIP, 15*DIP)}, MainStart, Stretch, 37 * DIP, 15 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, MainStart, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, MainCenter, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, MainEnd, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, SpaceAround, Stretch, 72 * DIP, 13 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, SpaceBetween, Stretch, 50 * DIP, 13 * DIP},
		{[]Element{mock(26*DIP, 13*DIP), mock(13*DIP, 11*DIP)}, Homogeneous, Stretch, (26*2 + 11) * DIP, 13 * DIP},
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
