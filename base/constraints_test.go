package base

import (
	"testing"
)

func TestConstraints(t *testing.T) {
	cases := []struct {
		in                                                   Constraints
		isNormalized, isTight, hasTightWidth, hasTightHeight bool
		isBounded, hasBoundedWidth, hasBoundedHeight         bool
	}{
		{Expand(), true, true, true, true, false, false, false},
		{ExpandHeight(10 * DIP), true, true, true, true, false, true, false},
		{ExpandWidth(10 * DIP), true, true, true, true, false, false, true},
		{Loose(Size{10 * DIP, 15 * DIP}), true, false, false, false, true, true, true},
		{Tight(Size{10 * DIP, 15 * DIP}), true, true, true, true, true, true, true},
		{TightWidth(10 * DIP), true, false, true, false, false, true, false},
		{TightHeight(10 * DIP), true, false, false, true, false, false, true},
	}

	for i, v := range cases {
		if out := v.in.IsNormalized(); v.isNormalized != out {
			t.Errorf("Failed on case %d for IsNormalized, want %v, got %v", i, v.isNormalized, out)
		}
		if out := v.in.IsTight(); v.isTight != out {
			t.Errorf("Failed on case %d for IsTight, want %v, got %v", i, v.isTight, out)
		}
		if out := v.in.HasTightWidth(); v.hasTightWidth != out {
			t.Errorf("Failed on case %d for HasTightWidth, want %v, got %v", i, v.hasTightWidth, out)
		}
		if out := v.in.HasTightHeight(); v.hasTightHeight != out {
			t.Errorf("Failed on case %d for HasTightHeight, want %v, got %v", i, v.hasTightHeight, out)
		}
		if out := v.in.IsBounded(); v.isBounded != out {
			t.Errorf("Failed on case %d for IsBounded, want %v, got %v", i, v.isBounded, out)
		}
		if out := v.in.HasBoundedWidth(); v.hasBoundedWidth != out {
			t.Errorf("Failed on case %d for HasBoundedWidth, want %v, got %v", i, v.hasBoundedWidth, out)
		}
		if out := v.in.HasBoundedHeight(); v.hasBoundedHeight != out {
			t.Errorf("Failed on case %d for HasBoundedHeight, want %v, got %v", i, v.hasBoundedHeight, out)
		}

		if out := v.in.IsZero(); out {
			t.Errorf("Failed on case %d for IsZero, want %v, got %v", i, false, out)
		}
	}
}

func TestConstraints_Constrain(t *testing.T) {
	cases := []struct {
		in   Constraints
		size Size
		out  Size
	}{
		{Tight(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(10 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{15 * DIP, 15 * DIP}},
		{TightWidth(5 * DIP), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 15 * DIP}},
		{TightHeight(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightHeight(30 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 30 * DIP}},
		{TightHeight(75 * DIP / 10), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 75 * DIP / 10}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{2 * DIP, 3 * DIP}},
	}

	for i, v := range cases {
		if out := v.in.Constrain(v.size); v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
		if out := v.in.ConstrainWidth(v.size.Width); v.out.Width != out {
			t.Errorf("Failed on case %d width, want %v, got %v", i, v.out.Width, out)
		}
		if out := v.in.ConstrainHeight(v.size.Height); v.out.Height != out {
			t.Errorf("Failed on case %d height, want %v, got %v", i, v.out.Height, out)
		}
	}
}

func TestConstraints_ConstrainAndAttemptToPreserveAspectRatio(t *testing.T) {
	cases := []struct {
		in   Constraints
		size Size
		out  Size
	}{
		{Tight(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(10 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightWidth(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{15 * DIP, 225 * DIP / 10}},
		{TightWidth(5 * DIP), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 75 * DIP / 10}},
		{TightHeight(15 * DIP), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{TightHeight(30 * DIP), Size{10 * DIP, 15 * DIP}, Size{20 * DIP, 30 * DIP}},
		{TightHeight(75 * DIP / 10), Size{10 * DIP, 15 * DIP}, Size{5 * DIP, 75 * DIP / 10}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 15 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{20 * DIP, 30 * DIP}, Size{10 * DIP, 15 * DIP}},
		{Loose(Size{10 * DIP, 15 * DIP}), Size{2 * DIP, 3 * DIP}, Size{2 * DIP, 3 * DIP}},
	}

	for i, v := range cases {
		out := v.in.ConstrainAndAttemptToPreserveAspectRatio(v.size)
		if v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}
func TestConstraints_Inset(t *testing.T) {
	cases := []struct {
		in      Constraints
		deflate Length
		out     Constraints
	}{
		{Tight(Size{}), 1 * DIP, Tight(Size{})},
		{Tight(Size{2 * DIP, 2 * DIP}), 10 * DIP, Tight(Size{})},
		{Tight(Size{10 * DIP, 11 * DIP}), 5 * DIP, Tight(Size{5 * DIP, 6 * DIP})},
		{Loose(Size{}), 1 * DIP, Loose(Size{})},
		{Loose(Size{2 * DIP, 2 * DIP}), 10 * DIP, Loose(Size{})},
		{Loose(Size{10 * DIP, 11 * DIP}), 5 * DIP, Loose(Size{5 * DIP, 6 * DIP})},
		{TightWidth(0), 1 * DIP, TightWidth(0)},
		{TightWidth(2 * DIP), 10 * DIP, TightWidth(0)},
		{TightWidth(10 * DIP), 5 * DIP, TightWidth(5 * DIP)},
		{TightHeight(0), 1 * DIP, TightHeight(0)},
		{TightHeight(2 * DIP), 10 * DIP, TightHeight(0)},
		{TightHeight(10 * DIP), 5 * DIP, TightHeight(5 * DIP)},
		{Expand(), 5 * DIP, Expand()},
	}

	for i, v := range cases {
		out := v.in.Inset(v.deflate, v.deflate)
		if v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}

func TestConstraints_Tighten(t *testing.T) {
	cases := []struct {
		in   Constraints
		size Size
		out  Constraints
	}{
		{Expand(), Size{10 * DIP, 10 * DIP}, Expand()},
		{ExpandHeight(10 * DIP), Size{10 * DIP, 10 * DIP}, Constraints{Size{10 * DIP, Inf}, Size{10 * DIP, Inf}}},
		{ExpandWidth(10 * DIP), Size{10 * DIP, 10 * DIP}, Constraints{Size{Inf, 10 * DIP}, Size{Inf, 10 * DIP}}},
		{Loose(Size{20 * DIP, 25 * DIP}), Size{10 * DIP, 10 * DIP}, Tight(Size{10 * DIP, 10 * DIP})},
		{Loose(Size{20 * DIP, 25 * DIP}), Size{30 * DIP, 30 * DIP}, Tight(Size{20 * DIP, 25 * DIP})},
		{Tight(Size{10 * DIP, 15 * DIP}), Size{10 * DIP, 10 * DIP}, Tight(Size{10 * DIP, 15 * DIP})},
		{TightWidth(15 * DIP), Size{10 * DIP, 10 * DIP}, Tight(Size{15 * DIP, 10 * DIP})},
		{TightHeight(15 * DIP), Size{10 * DIP, 10 * DIP}, Tight(Size{10 * DIP, 15 * DIP})},
	}

	for i, v := range cases {
		if out := v.in.Tighten(v.size); out != v.out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}
