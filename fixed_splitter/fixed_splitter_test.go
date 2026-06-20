package fixed_splitter

import (
	"testing"

	"github.com/1Vewton/text_splitter"
)

// Test whether the FixedSplitter implements TextSplitter
func TestInterface(t *testing.T) {
	var splitter interface{} = &FixedSplitter{}
	_, ok := splitter.(text_splitter.TextSplitter)
	if !ok {
		t.Errorf("FixedSplitter does not implements TextSplitter")
	}
}
