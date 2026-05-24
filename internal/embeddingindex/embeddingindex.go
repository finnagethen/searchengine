package embeddingindex

import (
	"encoding/binary"
	"fmt"
	"os"
	"slices"

	"github.com/finnagethen/searchengine/internal/utils"
)

// EmbeddingFileHeader represents the header of the binary file containing the embeddings.
type EmbeddingFileHeader struct {
	NumDocs   uint32
	Dimension uint32
}

type EmbeddingIndex struct {
	Header             EmbeddingFileHeader // header containing the number of documents and dimension of the embeddings
	DocumentEmbeddings [][]float32         // document index -> document embedding
}

// NewEmbeddingIndex returns a new EmbeddingIndex.
func NewEmbeddingIndex() *EmbeddingIndex {
	return &EmbeddingIndex{}
}

// LoadEmbeddings loads the document embeddings from a binary file containing the embeddings.
// The binary file should contain the following structure:
//
//	<num douments> (uint32)
//	<dimension> (uint32)
//	<embeddings> ([]float32)
func (index *EmbeddingIndex) LoadEmbeddings(embeddingsPath string) error {
	defer utils.Measure("EmbeddingIndex_LoadEmbeddings")()

	file, err := os.Open(embeddingsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = binary.Read(file, binary.LittleEndian, &index.Header); err != nil {
		return err
	}

	numDocs, dim := int(index.Header.NumDocs), int(index.Header.Dimension)

	embeddings := make([]float32, numDocs*dim)
	if err := binary.Read(file, binary.LittleEndian, embeddings); err != nil {
		return err
	}

	// Reshape the embeddings into [][]float32.
	vectors := make([][]float32, numDocs)
	for i := 0; i < numDocs; i++ {
		start := i * dim
		end := start + dim
		vectors[i] = embeddings[start:end]
	}
	index.DocumentEmbeddings = vectors

	return nil
}

// CosineSimilarity calculates the cosine similarity of a vector `v` of shape [N] and a matrix `m` of shape [M, N].
// Returns the cosine similarities as a vector of shape [M] or an error on failure.
// The vectors are expected to be normalized.
func (index *EmbeddingIndex) CosineSimilarity(v []float32, m [][]float32) ([]float32, error) {
	if len(v) != len(m[0]) {
		return nil, fmt.Errorf("vectors must have the same length")
	}
	if int(index.Header.Dimension) != len(v) {
		return nil, fmt.Errorf("vectors must have the same dimension as the embeddings")
	}

	result := make([]float32, len(m))
	for i, row := range m {
		var dotProduct float32
		for j := range v {
			dotProduct += v[j] * row[j]
		}
		result[i] = dotProduct
	}

	return result, nil
}

// TopKNeighbors returns the indices of the k most similar documents to the given document index.
// Assumes `k` is smaller than the number of documents.
func (index *EmbeddingIndex) TopKNeighbors(document int, k int) ([]int, error) {
	if len(index.DocumentEmbeddings) == 0 {
		return nil, fmt.Errorf("empty document embeddings")
	}

	cosineSimilarities, err := index.CosineSimilarity(index.DocumentEmbeddings[document], index.DocumentEmbeddings)
	if err != nil {
		return nil, err
	}

	// Create a slice of indices to save the permutation.
	indices := make([]int, len(cosineSimilarities))
	for i := range indices {
		indices[i] = i
	}

	// Sort the indices by cosine similarity in descending order.
	slices.SortFunc(indices, func(a, b int) int {
		if cosineSimilarities[a] > cosineSimilarities[b] {
			return -1
		} else if cosineSimilarities[a] < cosineSimilarities[b] {
			return 1
		}
		return 0
	})

	return indices[:k], nil
}

// norm returns the Euclidean length of a vector.
//func norm(v []float32) float32 {
//	var n float32
//	for _, x := range v {
//		n += x * x
//	}
//	if n == 0 {
//		return 0
//	}
//	return float32(math.Sqrt(float64(n)))
//}
