package embedding

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
	"golang.org/x/sync/errgroup"
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

// EmbedAll embeds all the text.
func EmbedAll(
	ctx context.Context,
	texts []string,
	client openai.Client,
	model string,
	dimension int,
) (*EmbedResult, error) {
	resultMap := NewEmbedResult()
	group, ctx := errgroup.WithContext(ctx)
	for i, text := range texts {
		tmpText := text
		idx := i
		group.Go(
			func() error {
				result, err := Embed(
					ctx,
					tmpText,
					client,
					model,
					dimension,
				)
				if err != nil {
					return err
				}
				resultMap.Set(
					idx,
					result,
				)
				return nil
			},
		)
	}
	if err := group.Wait(); err != nil {
		return nil, err
	}
	return resultMap, nil
}
