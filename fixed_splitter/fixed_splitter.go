// Split the document into chunks that have fixed sizes.
package fixed_splitter

import (
	"errors"

	"github.com/1Vewton/text_splitter"
)

// Fixed splitter data structure
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

// Create new fixed splitter
func New(
	chunk_size int,
	overlap int,
	document *string,
	documents *[]string,
) *FixedSplitter {
	return &FixedSplitter{
		ChunkSize: chunk_size,
		Overlap:   overlap,
		Document:  document,
		Documents: documents,
	}
}

// Split the single document
func (splitter *FixedSplitter) SplitText() ([]string, error) {
	var result []string = []string{}
	// Check if the document field is nil
	if splitter.Document == nil {
		return result, errors.New(
			"You cannot give a nil pointer to the Document field of splitter while using SplitText method",
		)
	}
	// If the real length is smaller than Chunksize
	runed_document := []rune(*splitter.Document)
	if len(runed_document) < splitter.ChunkSize {
		result = append(result, *splitter.Document)
		return result, nil
	}
	start := 0
	for start < len(runed_document) {
		end := max(start+splitter.ChunkSize, len(runed_document))
		result = append(result, string(runed_document[start:end]))
		start = start + splitter.ChunkSize - splitter.Overlap
	}
	// Default return if no error occurs
	return result, nil
}

// Split multiple documents
func (splitter *FixedSplitter) SplitMultipleTexts() ([]*text_splitter.SplitResult, error) {
	var result []*text_splitter.SplitResult = []*text_splitter.SplitResult{}
	// Return the default result
	return result, nil
}
