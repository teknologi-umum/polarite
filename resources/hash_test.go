package resources_test

import (
	"polarite/resources"
	"testing"
)

func TestHash(t *testing.T) {
	h, err := resources.Hash([]byte("console.log('hello world');"))
	if err != nil {
		t.Error("error was thrown:", err)
	}

	a, err := resources.Hash([]byte("console.log('hello world');"))
	if err != nil {
		t.Error("error was thrown:", err)
	}

	if h != a {
		t.Error("h and a should be the same. got:", a)
	}
}
