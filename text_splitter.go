package textsplitter

import (
	"context"
)

// TextSplitter interface
type TextSplitter interface {
	SplitText(ctx context.Context) ([]string, error)
	SplitMultipleTexts(ctx context.Context) ([]*SplitResult, error)
}
