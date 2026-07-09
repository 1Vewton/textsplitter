package semanticsplitter

import (
	"github.com/1Vewton/textsplitter/recursivesplitter"
)

// SemanticSplitter splits the text according to the cosine similarity of the meaning of two chunks
type SemanticSplitter struct {
	ChunkSize   int
	Overlap     int
	SubSplitter *recursivesplitter.RecursiveSplitter
}

// NewSemanticSplitter creates a new SemanticSplitter
func NewSemanticSplitter(
	ChunkSize int,
	Overlap int,
	SubSplitter *recursivesplitter.RecursiveSplitter,
) *SemanticSplitter {
	return &SemanticSplitter{
		ChunkSize:   ChunkSize,
		Overlap:     Overlap,
		SubSplitter: SubSplitter,
	}
}
