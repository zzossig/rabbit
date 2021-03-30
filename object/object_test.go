package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{"Hello World"}
	hello2 := &String{"Hello World"}
	diff1 := &String{"My name is johnny"}
	diff2 := &String{"My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{true}
	true2 := &Boolean{true}
	false1 := &Boolean{false}
	false2 := &Boolean{false}
	num1 := &Integer{1}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("trues should have same hash key")
	}

	if num1.HashKey() == true1.HashKey() {
		t.Errorf("num1 and true1 should not have same hash key")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("falses do not have same hash key")
	}

	if true1.HashKey() == false1.HashKey() {
		t.Errorf("true has same hash key as false")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{1}
	one2 := &Integer{1}
	two1 := &Integer{2}
	two2 := &Integer{2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}

func TestDecimalHashKey(t *testing.T) {
	one1 := &Decimal{1.111}
	one2 := &Decimal{1.111}
	two1 := &Decimal{2.222}
	two2 := &Decimal{2.222}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("integers with same content have twoerent hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("integers with twoerent content have same hash keys")
	}
}
