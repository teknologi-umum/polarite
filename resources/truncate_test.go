package resources_test

import (
	"polarite/resources"
	"testing"
)

func TestTruncateString(t *testing.T) {
	s := resources.TruncateString("Hello world", 5)
	if s != "Hello" {
		t.Error("s should be \"Hello\", got:", s)
	}
}

func TestTruncateString_Empty(t *testing.T) {
	s := resources.TruncateString("Hello", 0)
	if s != "" {
		t.Error("s should be empty, got:", s)
	}
}
