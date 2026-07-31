package semanticsplitter

import (
	"testing"

	"github.com/1Vewton/textsplitter"
)

// Test whether the SemanticSplitter implements TextSplitter
func TestInterface(t *testing.T) {
	var splitter interface{} = &SemanticSplitter{}
	_, ok := splitter.(textsplitter.TextSplitter)
	if !ok {
		t.Errorf("SemanticSplitter does not implements TextSplitter")
	}
}
