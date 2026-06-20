package fixed_splitter

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
