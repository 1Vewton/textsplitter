package vectorc

// VectorCOperator defines the operator for carrying out vector calculation using C
type CVectorOperator struct {
}

// NewCVectorOperator creates new CVectorOperator
func NewCVectorOperator() *CVectorOperator {
	return &CVectorOperator{}
}

// CosineSimularity carries out cosine simularity calculation
func (operator *CVectorOperator) CosineSimularity(
	vector1 []float64,
	vector2 []float64,
) (float64, error) {
	return CosineSimilarity(vector1, vector2)
}
