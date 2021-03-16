package mock

import "testing"

func TestConcatBits(t *testing.T) {
	in := []byte{
		0, 1, 1, 0, 1, 0, 1, 1,
	}
	var exp byte = 0b01101011
	result := concatBits(in)
	if result != exp {
		t.Errorf("expect: %v, got: %v", exp, result)
	}
}
