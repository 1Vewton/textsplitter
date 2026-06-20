// Package fixedsplitter splits the document into chunks that have fixed sizes.
package fixedsplitter

import (
	"context"
	"errors"

	"github.com/1Vewton/textsplitter"
)

// FixedSplitter split documents into chunks that have fixed sizes
type FixedSplitter struct {
	// Largest number of characters inside a chunk
	ChunkSize int
	// Maximum number of characters shared by two neighboring chunks
	Overlap int
	// Original Document
	Document *string
	// List of documents
	Documents *[]string
}

// NewFixedSplitter creates new fixed splitter
func NewFixedSplitter(
	chunkSize int,
	overlap int,
	document *string,
	documents *[]string,
) *FixedSplitter {
	return &FixedSplitter{
		ChunkSize: chunkSize,
		Overlap:   overlap,
		Document:  document,
		Documents: documents,
	}
}

// SplitText splits the single document
func (splitter *FixedSplitter) SplitText(ctx context.Context) ([]string, error) {
	var result []string = []string{}
	// Check if the document field is nil
	if splitter.Document == nil {
		return result, errors.New(
			"You cannot give a nil pointer to the Document field of splitter while using SplitText method",
		)
	}
	// If the real length is smaller than Chunksize
	runedDocument := []rune(*splitter.Document)
	if len(runedDocument) < splitter.ChunkSize {
		result = append(result, *splitter.Document)
		return result, nil
	}
	start := 0
	for start < len(runedDocument) {
		end := max(start+splitter.ChunkSize, len(runedDocument))
		result = append(result, string(runedDocument[start:end]))
		start = start + splitter.ChunkSize - splitter.Overlap
	}
	// Default return if no error occurs
	return result, nil
}

// SplitMultipleTexts splits multiple documents
func (splitter *FixedSplitter) SplitMultipleTexts(ctx context.Context) ([]*textsplitter.SplitResult, error) {
	var result []*textsplitter.SplitResult = []*textsplitter.SplitResult{}
	// Return the default result
	return result, nil
}
