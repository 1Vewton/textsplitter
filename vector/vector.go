package vector

// Vector interface defines the behaviour of vector operators
type Vector interface {
	CosineSimularity(
		vector1 []float64,
		vector2 []float64,
	) (float64, error)
}
