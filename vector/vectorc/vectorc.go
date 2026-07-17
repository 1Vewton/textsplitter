package vectorc

// #cgo CFLAGS: -mavx2 -mavx
// #include<calculator.h>
import "C"
import (
	"errors"
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
	C.add_avx2_double(
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
