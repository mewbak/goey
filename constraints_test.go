package goey

import (
	"testing"
)

func TestBox_Deflate(t *testing.T) {
	cases := []struct {
		in      Box
		deflate Length
		out     Box
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
		out := v.in.Deflate(v.deflate, v.deflate)
		if v.out != out {
			t.Errorf("Failed on case %d, want %v, got %v", i, v.out, out)
		}
	}
}
