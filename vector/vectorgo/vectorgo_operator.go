package vectorgo

// GoVectorOperator defines how the operator of vectorgo calculation method
type GoVectorOperator struct {
}

// NewGoVectorOperator creates new vector operator
func NewGoVectorOperator() *GoVectorOperator {
	return &GoVectorOperator{}
}

// CosineSimilarity calculates the cosine simularity of two vectors
func (operator *GoVectorOperator) CosineSimilarity(
	vector1 []float64,
	vector2 []float64,
) (float64, error) {
	return CosineSimilarity(vector1, vector2)
}
