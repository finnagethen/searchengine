package embeddingindex

import (
	"reflect"
	"testing"

	"github.com/finnagethen/searchengine/internal/utils"
)

func TestEmbeddingIndex_EmbedDocument(t *testing.T) {
	index := NewEmbeddingIndex()

	index.Vocab = map[string]int{"a": 0, "b": 1}
	// Flattened matrix of shape [2 x 2]:
	// [[1.0, 0.5], [-2.0, 1.5]]
	index.TokenEmbedings = []float64{
		1.0, 0.5, // a
		-2.0, 1.5, // b
	}
	index.Dimension = 2

	tests := []struct {
		name     string
		document string
		expected []float64
	}{
		{"empty_document", "", []float64{0.0, 0.0}},
		{"a_b_a", "a b a", []float64{0.0, 2.5}},
		{"b_b", "b b", []float64{-4.0, 3.0}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := index.EmbedDocument(tc.document)
			if err != nil {
				t.Fatalf("EmbedDocument failed: %v", err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("EmbedDocument(%q) = %+v; want %+v", tc.document, got, tc.expected)
			}
		})
	}
}

func TestEmbeddingIndex_CosineSimilarity(t *testing.T) {
	index := NewEmbeddingIndex()
	index.Dimension = 2

	// Flattened matrix of shape [4 x 2]:
	// [[10.0, 0.0], [0.0, 0.1], [8.0, 8.0], [-4.5, 0.0]]
	m := []float64{
		10.0, 0.0,
		0.0, 0.1,
		8.0, 8.0,
		-4.5, 0.0,
	}

	tests := []struct {
		name     string
		v        []float64
		expected []float64
	}{
		{"v1", []float64{4.0, 0.0}, []float64{1.000, 0.000, 0.707, -1.000}},
		{"v2", []float64{0.0, 0.0}, []float64{0.0, 0.0, 0.0, 0.0}},
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
	// Mock equivalent to {"a": [1.0, 0.5], "b": [-2.0, 1.5]}
	index.Vocab = map[string]int{"a": 0, "b": 1}
	index.TokenEmbedings = []float64{
		1.0, 0.5, // "a" (index 0)
		-2.0, 1.5, // "b" (index 1)
	}
	index.Dimension = 2

	documents := []string{"a a a b", "b b", "a b", ""}
	err := index.BuildFromDocuments(documents)
	if err != nil {
		t.Fatalf("BuildFromDocuments failed: %v", err)
	}

	tests := []struct {
		name     string
		document string
		k        int
		expected []int
	}{
		{"query_a_a_a_b", "a a a b", 2, []int{0, 2}},
		{"query_b", "b", 2, []int{1, 2}},
		{"query_b_top_0", "b", 0, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := index.TopKNeighbors(tc.document, tc.k)
			if err != nil {
				t.Fatalf("TopKNeighbors failed: %v", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("TopKNeighbors(%q, %d) = %+v; want %+v", tc.document, tc.k, got, tc.expected)
			}
		})
	}
}
