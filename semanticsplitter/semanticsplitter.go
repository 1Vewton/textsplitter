package semanticsplitter

import (
	"context"
	"errors"
	"fmt"

	"github.com/1Vewton/textsplitter/internal/embedding"
	"github.com/1Vewton/textsplitter/recursivesplitter"
	"github.com/1Vewton/textsplitter/vector"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// SemanticSplitter splits the text according to the cosine similarity of the meaning of two chunks
type SemanticSplitter struct {
	ChunkSize          int
	Overlap            int
	SubSplitter        *recursivesplitter.RecursiveSplitter
	EmbeddingClient    openai.Client
	EmbeddingModel     string
	EmbeddingDimension int
	VectorOperator     vector.Vector
	Threshold          float64
}

// NewSemanticSplitter creates a new SemanticSplitter
func NewSemanticSplitter(
	chunkSize int,
	overlap int,
	subSplitter *recursivesplitter.RecursiveSplitter,
	APIKey string,
	baseURL string,
	embeddingModel string,
	dimension int,
	vectorOperator vector.Vector,
	Threshold float64,
) *SemanticSplitter {
	if subSplitter.ChunkSize > chunkSize {
		fmt.Println(
			"[WARNING] It would be better to set the chunksize of the subSplitter to be smaller than the chunksize",
		)
	}
	return &SemanticSplitter{
		ChunkSize:   chunkSize,
		Overlap:     overlap,
		SubSplitter: subSplitter,
		EmbeddingClient: openai.NewClient(
			option.WithAPIKey(APIKey),
			option.WithBaseURL(baseURL),
		),
		EmbeddingModel:     embeddingModel,
		EmbeddingDimension: dimension,
		VectorOperator:     vectorOperator,
		Threshold:          Threshold,
	}
}

// embed embeds a piece of text
func (splitter *SemanticSplitter) embed(
	ctx context.Context,
	texts []string,
) (*embedding.EmbedResult, error) {
	return embedding.EmbedAll(
		ctx,
		texts,
		splitter.EmbeddingClient,
		splitter.EmbeddingModel,
		splitter.EmbeddingDimension,
	)
}

// SplitText splits a single text
func (splitter *SemanticSplitter) SplitText(
	ctx context.Context,
	document string,
) (
	[]string,
	error,
) {
	var result []string = []string{}
	// Split into sub chunks
	splitResult, errSplit := splitter.SubSplitter.SplitText(
		ctx,
		document,
	)
	if errSplit != nil {
		return nil, errSplit
	}
	if len(splitResult) < 1 {
		return nil, errors.New("The splitResult is empty")
	}
	// Embed the sub chunks
	embedResult, errEmbed := splitter.embed(
		ctx,
		splitResult,
	)
	if errEmbed != nil {
		return nil, errEmbed
	}
	// Start embedding
	tmpChunk := splitResult[0]
	for i := 1; i < len(splitResult); i++ {
		// Embed and similarity calculation
		result1, exists1 := embedResult.Read(i - 1)
		if !exists1 {
			return nil, fmt.Errorf(
				"Embed result %d does not exists",
				i-1,
			)
		}
		result2, exists2 := embedResult.Read(i)
		if !exists2 {
			return nil, fmt.Errorf(
				"Embed result %d does not exists",
				i,
			)
		}
		similarity, errCalc := splitter.VectorOperator.CosineSimilarity(
			result1,
			result2,
		)
		if errCalc != nil {
			return nil, errCalc
		}
		// Process
		chunkedTmpChunk := []rune(tmpChunk)
		chunkedDocument2 := []rune(splitResult[i])
		if similarity > splitter.Threshold &&
			len(chunkedDocument2)+len(chunkedTmpChunk) < splitter.ChunkSize {
			tmpChunk += splitResult[i]
		} else {
			result = append(result, tmpChunk)
			tmpChunk = splitResult[i]
		}
	}
	if tmpChunk != "" {
		result = append(result, tmpChunk)
	}
	return result, nil
}
