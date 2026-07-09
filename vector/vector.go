package vector

import (
	"errors"
)

// DotProduct calculates the dot product of two vectors
func DotProduct(
	vector1 []float64,
	vector2 []float64,
) (float64, error) {
	var result float64 = 0.0
	if len(vector1) != len(vector2) {
		return result, errors.New("The length of the vectors are not the same")
	}
	for i := 0; i < len(vector1); i++ {
		result += (vector1[i] * vector2[i])
	}
	return result, nil
}

// CosineSimularity calculates the cosine simularity of two vectors
func CosineSimularity(
	vector1 []float64,
	vector2 []float64,
) float64 {
	return 0.0
}
