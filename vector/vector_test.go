package vector

import (
	"testing"
)

// Test the length checking of dot product function
func TestDotProductLengthChecking(t *testing.T) {
	vector1 := []float64{0.1, 0.2, 0.3}
	vector2 := []float64{0.1, 0.2}
	_, err := DotProduct(vector1, vector2)
	if err == nil {
		t.Fatal("This function ought to through an error")
	}
}

// Test the result of dot product
func TestDotProduct(t *testing.T) {
	vector1 := []float64{1.0, 2.0, 3.0}
	result, err := DotProduct(vector1, vector1)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result != 14.0 {
		t.Errorf("The result calculated is not correct")
	}
}
