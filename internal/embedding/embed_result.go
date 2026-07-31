package embedding

import (
	"sync"
)

// EmbedResult stores the result of embedding
type EmbedResult struct {
	sync.RWMutex
	m map[int][]float64
}

// NewEmbedResult creates new EmbedResult
func NewEmbedResult() *EmbedResult {
	return &EmbedResult{
		m: make(map[int][]float64),
	}
}

// Set sets a key-value in EmbedResult
func (result *EmbedResult) Set(
	key int,
	value []float64,
) {
	result.Lock()
	defer result.Unlock()
	result.m[key] = value
}

// Read reads the value for key and see if it exists
func (result *EmbedResult) Read(
	key int,
) ([]float64, bool) {
	result.RLock()
	defer result.RUnlock()
	value, exists := result.m[key]
	return value, exists
}
