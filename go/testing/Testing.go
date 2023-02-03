package main

import "testing"

func TestAverage(t *testing.T) {
	v := Sum()
	if v != 1 {
		t.Error("Expected 1.5, got ", v)
	}
}

func Sum() int {
	return 1;
}