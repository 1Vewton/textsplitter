package semanticsplitter

import (
	"context"
	"fmt"

	"github.com/1Vewton/textsplitter/internal/embedding"
	"github.com/1Vewton/textsplitter/recursivesplitter"
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
	}
}

// embed embeds aa piece of text
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
