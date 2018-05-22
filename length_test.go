package goey

import (
	"testing"
)

func TestLength(t *testing.T) {
	if rt := (1 * Length).Length(); rt != 1 {
		t.Errorf("Unexpected round-trim for Length")
	}
	if rt := (1 * PT).PT(); rt != 1 {
		t.Errorf("Unexpected round-trim for PT")
	}
}
