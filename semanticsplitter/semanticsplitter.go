package semanticsplitter

import (
	"context"
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
	VectorOperator     *vector.Vector
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
	vectorOperator *vector.Vector,
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
	}
}

// embed embeds a piece of text
func (splitter *SemanticSplitter) embed(
	ctx context.Context,
	text string,
) ([]float64, error) {
	return embedding.Embed(
		ctx,
		text,
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
	_, errSplit := splitter.SubSplitter.SplitText(
		ctx,
		document,
	)
	if errSplit != nil {
		return nil, errSplit
	}
	return result, nil
}
