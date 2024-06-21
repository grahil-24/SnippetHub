package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// indicates to the Go test runner that our Equal() function is a test helper
	t.Helper()
	if actual != expected {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func StringContains(t *testing.T, actual, expected string) {
	t.Helper()

	if !strings.Contains(actual, expected) {
		t.Errorf("got: %q; expected: %q", actual, expected)
	}
}
