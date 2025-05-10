package utils

import "testing"

func Assert(t *testing.T, expected, received any) {
	if expected != received {
		t.Errorf("Expected: %s\nReceived: %s\n", expected, received)
	}
}
