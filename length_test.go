package goey

import (
	"testing"
)

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
