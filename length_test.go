package goey

import (
	"fmt"
	"testing"
)

func ExampleLength() {
	// Since there are 96 device-independent pixels per inch, and 6 picas
	// per inch, the following two lengths should be equal.
	length1 := 96 * DIP
	length2 := 6 * PC

	if length1 == length2 {
		fmt.Printf("All is OK with the world.")
	} else {
		fmt.Printf("This should not happen, unless there is a rounding error.")
	}

	// Output:
	// All is OK with the world.
}

func ExampleLength_String() {
	fmt.Printf("Converting:  1pt is equal to %sdip\n", 1*PT)
	fmt.Printf("Converting:  1pt is equal to %1.2fdip\n", (1 * PT).DIP())
	fmt.Printf("Converting:  1pc is equal to %1.1fdip\n", (1 * PC).DIP())

	// Output:
	// Converting:  1pt is equal to 1:21dip
	// Converting:  1pt is equal to 1.33dip
	// Converting:  1pc is equal to 16.0dip
}

func ExampleRectangle() {
	r := Rectangle{Point{10 * DIP, 20 * DIP}, Point{90 * DIP, 80 * DIP}}

	fmt.Printf("Rectangle %s has dimensions %.0fdip by %.0fdip.",
		r, r.Dx().DIP(), r.Dy().DIP(),
	)

	// Output:
	// Rectangle (10:00,20:00)-(90:00,80:00) has dimensions 80dip by 60dip.
}

func TestLength(t *testing.T) {
	t.Logf("Constants DIP: %d", 1*DIP)
	t.Logf("Constants PT:  %d", 1*PT)
	t.Logf("Constants PC:  %d", 1*PC)
	if rt := (1 * DIP).DIP(); rt != 1 {
		t.Errorf("Unexpected round-trim for Length, %v =/= %v", rt, 1)
	}
	if rt := (1 * PT).PT(); rt != 1 {
		t.Errorf("Unexpected round-trim for PT,  %v =/= %v", rt, 1)
	}
	if rt := (1 * PC).PC(); rt != 1 {
		t.Errorf("Unexpected round-trim for PT,  %v =/= %v", rt, 1)
	}
	if rt := (1 * PT) * (1 << 6) / (1 * DIP); rt != 96*(1<<6)/72 {
		t.Errorf("Unexpected ratio between DIP and PT, %v =/= %v", rt, 96*(1<<6)/72)
	}
	if rt := (1 * PC) * (1 << 6) / (1 * DIP); rt != 96*(1<<6)/6 {
		t.Errorf("Unexpected ratio between DIP and PC, %v =/= %v", rt, 96*(1<<6)/72)
	}
}

func TestLength_Clamp(t *testing.T) {
	cases := []struct {
		in       Length
		min, max Length
		out      Length
	}{
		{10 * DIP, 0 * DIP, 20 * DIP, 10 * DIP},
		{30 * DIP, 0 * DIP, 20 * DIP, 20 * DIP},
		{-10 * DIP, 0 * DIP, 20 * DIP, 0 * DIP},
		{10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{30 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{-10 * DIP, 10 * DIP, 10 * DIP, 10 * DIP},
		{10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{30 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
		{-10 * DIP, 20 * DIP, 0 * DIP, 20 * DIP},
	}

	for i, v := range cases {
		if out := v.in.Clamp(v.min, v.max); out != v.out {
			t.Errorf("Error in case %d, want %s, got %s", i, v.out, out)
		}
	}
}
