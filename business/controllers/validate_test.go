package controllers_test

import (
	"math/rand"
	"polarite/business/controllers"
	"testing"
)

func TestValidateSize(t *testing.T) {
	a := controllers.ValidateSize([]byte("Hello world"))
	if a != false {
		t.Error("should be false, got:", a)
	}

	z := make([]byte, 1024*6)
	rand.Read(z)
	b := controllers.ValidateSize(z)
	if b != true {
		t.Error("should be true, got:", b)
	}
}

func TestValidateID(t *testing.T) {
	arr := []string{"some", "random", "text"}

	a := controllers.ValidateID(arr, "text")
	if a != true {
		t.Error("should be true, got:", a)
	}

	b := controllers.ValidateID(arr, "blah")
	if b != false {
		t.Error("should be false, got:", b)
	}
}
