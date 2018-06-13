package goey

import (
	"fmt"
	"testing"
)

func ExampleLength() {
	// Since there are 96 device-independent pixels per inch, and 72 points
	// per inch, the following two lengths should be equal.
	length1 := 96 * DIP
	length2 := 72 * DIP

	if length1 == length2 {
		fmt.Printf("All is OK with the world.")
	} else {
		fmt.Printf("This should not happen, unless there is a bad rounding error.")
	}
}

func ExampleLength_String() {
	fmt.Printf("Converting:  1pt is equal to %sdip\n", 1*PT)
	fmt.Printf("Converting:  1pt is equal to %1.2fdip\n", (1 * PT).DIP())

	// Output:
	// Converting:  1pt is equal to 1:21dip
	// Converting:  1pt is equal to 1.33dip
}

func TestLength(t *testing.T) {
	t.Logf("Constants DIP: %d", 1*DIP)
	t.Logf("Constants PT:  %d", 1*PT)
	if rt := (1 * DIP).DIP(); rt != 1 {
		t.Errorf("Unexpected round-trim for Length, %v =/= %v", rt, 1)
	}
	if rt := (1 * PT).PT(); rt != 1 {
		t.Errorf("Unexpected round-trim for PT,  %v =/= %v", rt, 1)
	}
	if rt := (1 * PT) * (1 << 6) / (1 * DIP); rt != 96*(1<<6)/72 {
		t.Errorf("Unexpected ratio between DIP and PT, %v =/= %v", rt, 96*(1<<6)/72)
	}
}
