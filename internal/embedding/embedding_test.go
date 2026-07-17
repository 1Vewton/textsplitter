package embedding

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// Test the embedding function
func TestEmbedding(t *testing.T) {
	data, err := os.ReadFile(
		"testdata/response.json",
	)
	if err != nil {
		t.Fatalf(
			"The test failed when reading file due to %s",
			err.Error(),
		)
	}
	mockServer := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write(data)
			},
		),
	)
	defer mockServer.Close()
	var result []float64
	result, err = Embed(
		t.Context(),
		"The food was delicious and the waiter...",
		openai.NewClient(
			option.WithBaseURL(mockServer.URL),
			option.WithAPIKey("114514"),
		),
		"text-embedding-ada-002",
		3,
	)
	if err != nil {
		t.Error(err.Error())
	}
	if result[0] != 0.0023064255 &&
		result[1] != -0.009327292 &&
		result[2] != -0.0028842222 {
		t.Errorf("There is a problem with output")
	}
	if len(result) != 3 {
		t.Errorf("Expected length 3. got %d", len(result))
	}
}
