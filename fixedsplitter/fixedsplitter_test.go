package fixedsplitter

import (
	"testing"

	"github.com/1Vewton/textsplitter"
)

// Test whether the FixedSplitter implements TextSplitter
func TestInterface(t *testing.T) {
	var splitter interface{} = &FixedSplitter{}
	_, ok := splitter.(textsplitter.TextSplitter)
	if !ok {
		t.Errorf("FixedSplitter does not implements TextSplitter")
	}
}
