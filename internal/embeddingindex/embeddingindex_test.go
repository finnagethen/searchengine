package embeddingindex

import (
	"reflect"
	"testing"

	"github.com/finnagethen/searchengine/internal/utils"
)

func TestEmbeddingIndex_CosineSimilarity(t *testing.T) {
	index := NewEmbeddingIndex()
	index.Header.Dimension = 2

	// Matrix of shape [4 x 2].
	// The vectors are pre-normalized so the dot product accurately reflects cosine similarity.
	m := [][]float32{
		{1.0, 0.0},               // [1, 0]
		{0.0, 1.0},               // [0, 1]
		{0.70710678, 0.70710678}, // [1/sqrt(2), 1/sqrt(2)]
		{-1.0, 0.0},              // [-1, 0]
	}

	tests := []struct {
		name     string
		v        []float32
		expected []float32
	}{
		{"v1_normalized", []float32{1.0, 0.0}, []float32{1.000, 0.000, 0.707, -1.000}},
		{"v2_zeros", []float32{0.0, 0.0}, []float32{0.0, 0.0, 0.0, 0.0}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := index.CosineSimilarity(tc.v, m)
			if err != nil {
				t.Fatalf("CosineSimilarity failed: %v", err)
			}

			if !utils.EqualSliceEpsilon(got, tc.expected, 1e-3) {
				t.Errorf("CosineSimilarity(%+v, m) = %+v; want %+v", tc.v, got, tc.expected)
			}
		})
	}
}

func TestEmbeddingIndex_TopKNeighbors(t *testing.T) {
	index := NewEmbeddingIndex()
	index.Header.Dimension = 2

	// Pre-populate DocumentEmbeddings to mock a loaded binary file.
	// Using normalized vectors for accurate similarity sorting.
	index.DocumentEmbeddings = [][]float32{
		{1.0, 0.0},      // Doc 0 (Target)
		{0.0, 1.0},      // Doc 1 (Orthogonal to 0)
		{0.995, 0.0998}, // Doc 2 (Highly similar to 0)
		{-1.0, 0.0},     // Doc 3 (Exact opposite of 0)
	}

	tests := []struct {
		name     string
		document int
		k        int
		expected []int
	}{
		{"query_doc_0_top_2", 0, 2, []int{0, 2}},
		{"query_doc_1_top_2", 1, 2, []int{1, 2}},
		{"query_doc_0_top_0", 0, 0, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := index.TopKNeighbors(tc.document, tc.k)
			if err != nil {
				t.Fatalf("TopKNeighbors failed: %v", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("TopKNeighbors(%d, %d) = %+v; want %+v", tc.document, tc.k, got, tc.expected)
			}
		})
	}
}
