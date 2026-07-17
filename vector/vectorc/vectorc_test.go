package vectorc

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
