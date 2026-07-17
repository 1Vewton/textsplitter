package vectorc

// #cgo CFLAGS: -mavx2 -mavx
// #include<calculator.h>
import "C"
import (
	"errors"
	"math"
	"unsafe"
)

// DotProduct calculates the dot product of two vectors using C
func DotProduct(
	vector1 []float64,
	vector2 []float64,
) (float64, error) {
	if len(vector1) != len(vector2) {
		return 0.0, errors.New("The length of the vectors are not the same")
	}
	out := make([]float64, len(vector1))
	C.mult_avx2_double(
		(*C.double)(unsafe.Pointer(&vector1[0])),
		(*C.double)(unsafe.Pointer(&vector2[0])),
		(*C.double)(unsafe.Pointer(&out[0])),
		C.int(len(vector1)),
	)
	result := C.sum_avx2_double(
		(*C.double)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return float64(result), nil
}

// EuclideanNorm calculates the euclidean norm of a vector using C
func EuclideanNorm(
	vector []float64,
) float64 {
	out := make([]float64, len(vector))
	C.mult_avx2_double(
		(*C.double)(unsafe.Pointer(&vector[0])),
		(*C.double)(unsafe.Pointer(&vector[0])),
		(*C.double)(unsafe.Pointer(&out[0])),
		C.int(len(vector)),
	)
	result := C.sum_avx2_double(
		(*C.double)(unsafe.Pointer(&out[0])),
		C.int(len(out)),
	)
	return math.Sqrt(float64(result))
}

// CosineSimularity calculates the cosine simularity of two vectors
func CosineSimularity(
	vector1 []float64,
	vector2 []float64,
) (float64, error) {
	dotProductResult, err := DotProduct(
		vector1,
		vector2,
	)
	if err != nil {
		return 0.0, err
	}
	under := EuclideanNorm(vector1) * EuclideanNorm(vector2)
	if under == 0.0 {
		return 0.0, errors.New(
			"The result of multiplying the euclidean norm of two vectors cannot be 0",
		)
	}
	return dotProductResult / under, nil
}
