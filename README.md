# text_splitter

**text_splitter** is a Go package for splitting documents into smaller chunks, designed for use in Retrieval-Augmented Generation (RAG) pipelines and other text-processing workflows.

The package defines a generic `TextSplitter` interface and provides two concrete implementations:

- **`FixedSplitter`** -- splits text into fixed-size chunks with configurable overlap.
- **`RecursiveSplitter`** -- splits text recursively along natural separators (paragraphs, sentences, words, etc.) to preserve semantic boundaries, with fixed-size chunking as a fallback.

## Features

- **Interface-driven design** -- Implement your own splitting strategy by satisfying the `TextSplitter` interface.
- **Fixed-size chunking** -- Split documents into chunks of a configurable maximum length.
- **Recursive splitting** -- Split text along natural language separators (e.g., `\n\n`, `\n`, `。`, `，`, ` `, `,`, `.`) to keep semantically related text together.
- **Configurable overlap** -- Share a specified number of characters between neighboring chunks to preserve context.
- **Concurrent multi-document splitting** -- Split multiple documents in parallel using `errgroup`, with automatic error propagation and context cancellation.
- **Unicode-aware** -- Chunk boundaries are calculated on rune count, so multi-byte characters (Chinese, emoji, etc.) are handled correctly.
- **Context support** -- All splitting methods accept `context.Context` for timeout and cancellation control.

## Installation

```bash
go get github.com/1Vewton/textsplitter
```

## Quick Start

### FixedSplitter -- Splitting a single document

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/fixedsplitter"
)

func main() {
    splitter := fixedsplitter.NewFixedSplitter(
        100, // ChunkSize  -- max characters per chunk
        20,  // Overlap    -- characters shared between adjacent chunks
    )

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    document := "This is a long document that needs to be split into smaller chunks for processing."

    chunks, err := splitter.SplitText(ctx, document)
    if err != nil {
        panic(err)
    }

    for i, chunk := range chunks {
        fmt.Printf("Chunk %d: %s\n", i+1, chunk)
    }
}
```

### FixedSplitter -- Splitting multiple documents (concurrently)

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/fixedsplitter"
)

func main() {
    splitter := fixedsplitter.NewFixedSplitter(100, 20)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    documents := []string{
        "First document content...",
        "Second document content...",
        "Third document content...",
    }

    results, err := splitter.SplitMultipleTexts(ctx, documents)
    if err != nil {
        panic(err)
    }

    for _, result := range results {
        fmt.Printf("FullText: %s\n", result.FullText)
        fmt.Printf("Chunk:    %s\n\n", result.ChunkResult)
    }
}
```

### RecursiveSplitter -- Splitting a single document

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/recursivesplitter"
)

func main() {
    splitter := recursivesplitter.NewRecursiveSplitter(
        100,                        // ChunkSize  -- max characters per chunk
        20,                         // Overlap    -- characters shared between adjacent chunks
        []string{"\n\n", "\n", "。", "，", " ", ",", "."}, // Separators -- split in order of priority
    )

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    document := "This is a long document that needs to be split into smaller chunks for processing."

    chunks, err := splitter.SplitText(ctx, document)
    if err != nil {
        panic(err)
    }

    for i, chunk := range chunks {
        fmt.Printf("Chunk %d: %s\n", i+1, chunk)
    }
}
```

### RecursiveSplitter -- Splitting multiple documents (concurrently)

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/recursivesplitter"
)

func main() {
    splitter := recursivesplitter.NewRecursiveSplitter(100, 20,
        []string{"\n\n", "\n", "。", "，", " ", ",", "."},
    )

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    documents := []string{
        "First document content...",
        "Second document content...",
        "Third document content...",
    }

    results, err := splitter.SplitMultipleTexts(ctx, documents)
    if err != nil {
        panic(err)
    }

    for _, result := range results {
        fmt.Printf("FullText: %s\n", result.FullText)
        fmt.Printf("Chunk:    %s\n\n", result.ChunkResult)
    }
}
```

## API

### Interface: `TextSplitter`

Defined in the root `textsplitter` package:

```go
type TextSplitter interface {
    SplitText(ctx context.Context, document string) ([]string, error)
    SplitMultipleTexts(ctx context.Context, documents []string) ([]*SplitResult, error)
}
```

| Method                 | Description                                          |
|------------------------|------------------------------------------------------|
| `SplitText`            | Splits a single document into a slice of chunk strings. |
| `SplitMultipleTexts`   | Splits multiple documents concurrently and returns a flat slice of `SplitResult` (one entry per chunk, each referencing its original full text). |

### Struct: `SplitResult`

```go
type SplitResult struct {
    FullText    string   // The original document text
    ChunkResult string   // One chunk from that document
}
```

### Struct: `FixedSplitter`

```go
type FixedSplitter struct {
    ChunkSize int     // Maximum number of characters in each chunk
    Overlap   int     // Number of characters shared between consecutive chunks
}
```

#### Constructor

```go
func NewFixedSplitter(chunkSize int, overlap int) *FixedSplitter
```

#### Methods

```go
func (splitter *FixedSplitter) SplitText(ctx context.Context, document string) ([]string, error)
func (splitter *FixedSplitter) SplitMultipleTexts(ctx context.Context, documents []string) ([]*textsplitter.SplitResult, error)
```

### Struct: `RecursiveSplitter`

```go
type RecursiveSplitter struct {
    ChunkSize  int        // Maximum number of characters in each chunk
    Overlap    int        // Number of characters shared between consecutive chunks
    Separators []string   // Ordered list of separators to split on (e.g., ["\n\n", "\n", "。", "，", " ", ",", "."])
}
```

#### Constructor

```go
func NewRecursiveSplitter(chunkSize int, overlap int, separators []string) *RecursiveSplitter
```

#### Methods

```go
func (splitter *RecursiveSplitter) SplitText(ctx context.Context, document string) ([]string, error)
func (splitter *RecursiveSplitter) SplitMultipleTexts(ctx context.Context, documents []string) ([]*textsplitter.SplitResult, error)
```

### Splitting Behavior

#### FixedSplitter

- If the document length (in runes) is **less than or equal to `ChunkSize`**, the entire document is returned as a single chunk.
- Otherwise, the document is split into chunks of at most `ChunkSize` runes.
- Consecutive chunks overlap by `Overlap` runes, preserving context across chunk boundaries.

Example with `ChunkSize=60`, `Overlap=20`:

```
Chunk 1: [characters 0-60)
Chunk 2: [characters 40-100)
Chunk 3: [characters 80-140)
...
```

#### RecursiveSplitter

- If the document length (in runes) is **less than or equal to `ChunkSize`**, the entire document is returned as a single chunk.
- Otherwise, the text is recursively split using the provided `Separators` list **in order of priority**:
  1. Try the first separator (e.g., `\n\n` for paragraphs). If the resulting parts fit within `ChunkSize`, they are merged up to the limit.
  2. If a part is still too long, move to the next separator (e.g., `\n` for lines) and recurse.
  3. If all separators are exhausted, fall back to `FixedSplitter` for a forced fixed-size split.
- This approach keeps semantically related text together as much as possible before resorting to character-level splitting.

## Concurrent Multi-Document Processing

`SplitMultipleTexts` uses [`errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup) to split each document concurrently. Benefits:

- **Parallel execution** across all documents.
- **Context cancellation** -- if one split fails or the context expires, all goroutines are cancelled.
- **Error propagation** -- the first non-nil error is returned.

## Unicode Support

The implementation internally converts strings to `[]rune` before chunking, ensuring that multi-byte characters (e.g., Chinese, Japanese, emoji) are counted correctly. Chunk boundaries never split a character's byte sequence.

## Requirements

- Go 1.25 or later

## License

[MIT](LICENSE)