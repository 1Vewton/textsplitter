package recursivesplitter

import (
	"context"
)

// RecursiveSplitter splits the text according to the natural order of the language.
type RecursiveSplitter struct {
	ChunkSize  int
	Overlap    int
	Separators []string
}

// NewRecursiveSplitter creates new recursive splitter
func NewRecursiveSplitter(
	chunkSize int,
	overlap int,
	separators []string,
) *RecursiveSplitter {
	return &RecursiveSplitter{
		ChunkSize:  chunkSize,
		Overlap:    overlap,
		Separators: separators,
	}
}

// recursiveSplit splits single document using recursive splitting method for each step
func (splitter *RecursiveSplitter) recursiveSplit(
	document string,
	sepIdx int,
) (
	[]string,
	error,
) {
	var result []string = []string{}
	runedDocument := []rune(document)
	// Directly return the document if the length of the document is smaller than ChunkSize
	if len(runedDocument) <= splitter.ChunkSize {
		result = append(result, document)
		return result, nil
	}
	return result, nil
}

// SplitText splits single document using recursive splitting method
func (splitter *RecursiveSplitter) SplitText(
	ctx context.Context,
	document string,
) (
	[]string,
	error,
) {
	result, err := splitter.recursiveSplit(document, 0)
	return result, err
}
