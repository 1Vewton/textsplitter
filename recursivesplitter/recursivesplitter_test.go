package recursivesplitter

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"errors"

	"github.com/1Vewton/textsplitter"
	"golang.org/x/sync/errgroup"
)

// Test whether the RecursiveSplitter implements TextSplitter
func TestInterface(t *testing.T) {
	var splitter interface{} = &RecursiveSplitter{}
	_, ok := splitter.(textsplitter.TextSplitter)
	if !ok {
		t.Errorf("RecursiveSplitter does not implements TextSplitter")
	}
}

// Test the split text method
func TestSplitText(t *testing.T) {
	chunkSize := 60
	overlap := 20
	content, errRead := os.ReadFile("testdata/split_text_1.md")
	if errRead != nil {
		t.Fatalf("Fatal error occured when reading testdata due to %s", errRead)
	}
	document := string(content)
	splitter := NewRecursiveSplitter(
		chunkSize,
		overlap,
		[]string{"\n\n", "\n", "。", "，", " ", ",", "."},
	)
	// Timeout checking
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	result, errChunk := splitter.SplitText(
		ctx,
		document,
	)
	if errChunk != nil {
		t.Fatalf("Fatal error occured when running test due to %s", errChunk)
	}
	t.Log(result)
	t.Log(len(result))
	for _, i := range result {
		if !strings.Contains(document, i) {
			t.Errorf("%s does not exists in original document", i)
		}
	}
}

// Test the SplitMultipleTexts
func TestSplitMultipleTexts(t *testing.T) {
	chunkSize := 60
	overlap := 20
	testTasks := []string{
		"testdata/split_text_1.md",
		"testdata/split_text_2.md",
		"testdata/split_text_3.md",
		"testdata/split_text_4.md",
	}
	contentChannel := make(chan []byte, len(testTasks))
	group, _ := errgroup.WithContext(t.Context())
	for _, taskFile := range testTasks {
		currentTaskFile := taskFile
		group.Go(
			func() error {
				content, err := os.ReadFile(currentTaskFile)
				if err == nil {
					select {
					case contentChannel <- content:
					default:
						return errors.New(
							"There is a problem with the length of contentChannel",
						)
					}
				}
				return err
			},
		)
	}
	if err := group.Wait(); err != nil {
		t.Fatalf("Testing failed due to %s", err)
	}
	var taskContents []string = []string{}
	close(contentChannel)
	for byteContent := range contentChannel {
		taskContents = append(taskContents, string(byteContent))
	}
	splitter := NewRecursiveSplitter(
		chunkSize,
		overlap,
		[]string{"\n\n", "\n", "。", "，", " ", ",", "."},
	)
	// Timeout checking
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	result, errChunk := splitter.SplitMultipleTexts(
		ctx,
		taskContents,
	)
	if errChunk != nil {
		t.Fatalf("Fatal error occured when running test due to %s", errChunk)
	}
	for _, i := range result {
		t.Log(i.ChunkResult)
		if !strings.Contains(i.FullText, i.ChunkResult) {
			t.Errorf("%s does not exists", i)
		}
	}
}
