# text_splitter

**text_splitter** is a Go package for splitting documents into smaller chunks, designed for use in Retrieval-Augmented Generation (RAG) pipelines and other text-processing workflows.

The package defines a generic `TextSplitter` interface and provides three concrete implementations:

- **`FixedSplitter`** -- splits text into fixed-size chunks with configurable overlap.
- **`RecursiveSplitter`** -- splits text recursively along natural separators (paragraphs, sentences, words, etc.) to preserve semantic boundaries, with fixed-size chunking as a fallback.
- **`SemanticSplitter`** -- splits text by embedding sub-chunks and merging adjacent ones based on the cosine similarity of their meaning.

## Features

- **Interface-driven design** -- Implement your own splitting strategy by satisfying the `TextSplitter` interface.
- **Fixed-size chunking** -- Split documents into chunks of a configurable maximum length.
- **Recursive splitting** -- Split text along natural language separators (e.g., `\n\n`, `\n`, `。`, `，`, ` `, `,`, `.`) to keep semantically related text together.
- **Semantic splitting** -- Embed sub-chunks and merge adjacent ones based on cosine similarity of their meaning, so semantically related content stays together.
- **Configurable overlap** -- Share a specified number of characters between neighboring chunks to preserve context.
- **Concurrent multi-document splitting** -- Split multiple documents in parallel using `errgroup`, with automatic error propagation and context cancellation.
- **Unicode-aware** -- Chunk boundaries are calculated on rune count, so multi-byte characters (Chinese, emoji, etc.) are handled correctly.
- **Context support** -- All splitting methods accept `context.Context` for timeout and cancellation control.
- **Pluggable vector operators** -- Choose between a pure-Go implementation (`vectorgo`) or a C implementation (`vectorc`) for cosine similarity calculation.

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

### SemanticSplitter -- Splitting a single document

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/recursivesplitter"
    "github.com/1Vewton/textsplitter/semanticsplitter"
    "github.com/1Vewton/textsplitter/vector/vectorgo"
)

func main() {
    // First, create the sub-splitter that produces the initial candidate chunks.
    subSplitter := recursivesplitter.NewRecursiveSplitter(
        100,
        20,
        []string{"\n\n", "\n", "。", "，", " ", ",", "."},
    )

    splitter := semanticsplitter.NewSemanticSplitter(
        300,                // ChunkSize         -- max runes for a merged semantic chunk
        0,                  // Overlap           -- runes shared between adjacent chunks
        subSplitter,        // SubSplitter       -- produces the initial sub-chunks
        "your-api-key",     // APIKey            -- OpenAI-compatible API key
        "https://api.openai.com/v1", // BaseURL -- OpenAI-compatible API base URL
        "text-embedding-3-small",    // EmbeddingModel    -- embedding model name
        1536,                        // Dimension         -- embedding vector dimension
        vectorgo.NewGoVectorOperator(), // VectorOperator -- cosine similarity implementation
        0.5, // Threshold -- merge sub-chunks when similarity is above this value
    )

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    document := "This is a long document that needs to be split " +
        "into semantically meaningful chunks."

    chunks, err := splitter.SplitText(ctx, document)
    if err != nil {
        panic(err)
    }

    for i, chunk := range chunks {
        fmt.Printf("Chunk %d: %s\n", i+1, chunk)
    }
}
```

### SemanticSplitter -- Splitting multiple documents (concurrently)

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/1Vewton/textsplitter/recursivesplitter"
    "github.com/1Vewton/textsplitter/semanticsplitter"
    "github.com/1Vewton/textsplitter/vector/vectorgo"
)

func main() {
    subSplitter := recursivesplitter.NewRecursiveSplitter(
        100,
        20,
        []string{"\n\n", "\n", "。", "，", " ", ",", "."},
    )

    splitter := semanticsplitter.NewSemanticSplitter(
        300,
        0,
        subSplitter,
        "your-api-key",
        "https://api.openai.com/v1",
        "text-embedding-3-small",
        1536,
        vectorgo.NewGoVectorOperator(),
        0.5,
    )

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

### Struct: `SemanticSplitter`

```go
type SemanticSplitter struct {
    ChunkSize          int                                // Maximum number of characters in each final chunk
    Overlap            int                                // Number of characters shared between consecutive chunks
    SubSplitter        *recursivesplitter.RecursiveSplitter // Splits the document into initial sub-chunks
    EmbeddingClient    openai.Client                      // OpenAI-compatible embedding client
    EmbeddingModel     string                             // Embedding model name (e.g., "text-embedding-3-small")
    EmbeddingDimension int                                // Dimension of the embedding vectors
    VectorOperator     vector.Vector                      // Cosine similarity implementation
    Threshold          float64                            // Merge threshold; adjacent sub-chunks with similarity above this are merged
}
```

#### Constructor

```go
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
) *SemanticSplitter
```

#### Methods

```go
func (splitter *SemanticSplitter) SplitText(ctx context.Context, document string) ([]string, error)
func (splitter *SemanticSplitter) SplitMultipleTexts(ctx context.Context, documents []string) ([]*textsplitter.SplitResult, error)
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

#### SemanticSplitter

- The document is first split into **sub-chunks** by the configured `SubSplitter` (typically a `RecursiveSplitter`).
- Each sub-chunk is **embedded** into a vector using the configured embedding model (OpenAI-compatible API).
- Adjacent sub-chunks are compared using **cosine similarity** (via the configured `VectorOperator`):
  - If the similarity is **above `Threshold`** and the combined sub-chunks fit within `ChunkSize` runes, they are **merged** into one chunk.
  - Otherwise, the current chunk is finalized and a new chunk starts.
- The merging only checks character-length limits (rune count); the `Overlap` field is reserved for future use and does not currently affect the merge logic.
- This produces chunks whose boundaries follow semantic meaning rather than fixed character positions.

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