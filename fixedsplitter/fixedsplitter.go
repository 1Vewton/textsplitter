// Package fixedsplitter splits the document into chunks that have fixed sizes.
package fixedsplitter

import (
	"context"
	"errors"

	"github.com/1Vewton/textsplitter"
	"golang.org/x/sync/errgroup"
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
) *FixedSplitter {
	return &FixedSplitter{
		ChunkSize: chunkSize,
		Overlap:   overlap,
	}
}

// SplitText splits the single document using FixedSplitter method
func (splitter *FixedSplitter) SplitText(
	ctx context.Context,
	document string,
) (
	[]string,
	error,
) {
	var result []string = []string{}
	// If the real length is smaller than Chunksize
	runedDocument := []rune(document)
	if len(runedDocument) < splitter.ChunkSize {
		result = append(result, document)
		return result, nil
	}
	start := 0
	for start < len(runedDocument) {
		end := min(start+splitter.ChunkSize, len(runedDocument))
		result = append(result, string(runedDocument[start:end]))
		start = start + splitter.ChunkSize - splitter.Overlap
	}
	// Default return if no error occurs
	return result, nil
}

// SplitMultipleTexts splits multiple documents using FixedSplitter method
func (splitter *FixedSplitter) SplitMultipleTexts(
	ctx context.Context,
	documents []string,
) (
	[]*textsplitter.SplitResult,
	error,
) {
	var result []*textsplitter.SplitResult = []*textsplitter.SplitResult{}
	resultChannel := make(
		chan *textsplitter.TempSplitResult,
		len(documents)*2,
	)
	// Executing the error group that split each document in documents
	group, ctx := errgroup.WithContext(ctx)
	for _, fullText := range documents {
		fullTextTmp := fullText
		group.Go(
			func() error {
				result, err := splitter.SplitText(
					ctx,
					fullTextTmp,
				)
				if err == nil {
					tmpSplitResult := &textsplitter.TempSplitResult{
						FullText:    fullTextTmp,
						ChunkResult: result,
					}
					select {
					case resultChannel <- tmpSplitResult:
					default:
						return errors.New(
							"There is a problem with length of resultChannel",
						)
					}
				}
				return err
			},
		)
	}
	if err := group.Wait(); err != nil {
		return result, err
	}
	// The channel has to be closed to process data
	close(resultChannel)
	// Process return data
	for chunkResult := range resultChannel {
		for _, chunk := range chunkResult.ChunkResult {
			tmpResult := &textsplitter.SplitResult{
				FullText:    chunkResult.FullText,
				ChunkResult: chunk,
			}
			result = append(result, tmpResult)
		}
	}
	// Return the default result
	return result, nil
}
