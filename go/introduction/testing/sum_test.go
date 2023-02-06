package testing

import "testing"

/**
Para testear todos los archivos deben acabat en _test
y para ejecutarlo go test
 */
func TestSum(t *testing.T) {
	v := Sum()
	if v != 1 {
		t.Error("Expected 1.5, got ", v)
	}
}

func Sum() int {
	return 1
}