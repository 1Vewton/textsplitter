package embedding

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
)

// Embed embeds aa piece of text
func Embed(
	ctx context.Context,
	text string,
	client openai.Client,
	model string,
	dimension int,
) ([]float64, error) {
	result, err := client.Embeddings.New(
		ctx,
		openai.EmbeddingNewParams{
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(text),
			},
			Model:          model,
			Dimensions:     openai.Int(int64(dimension)),
			EncodingFormat: openai.EmbeddingNewParamsEncodingFormatFloat,
		},
	)
	if err != nil {
		return nil, err
	}
	if len(result.Data) < 1 {
		return nil, errors.New("The length of the return data is smaller than 1")
	}
	return result.Data[0].Embedding, nil
}
