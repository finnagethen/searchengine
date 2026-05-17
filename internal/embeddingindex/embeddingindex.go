package embeddingindex

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"

	"github.com/finnagethen/searchengine/internal/utils"
)

type EmbeddingIndex struct {
	Vocab              map[string]int // token -> token index
	TokenEmbedings     []float64      // token index -> embedding
	DocumentEmbeddings []float64      // document index -> embedding
	Dimension          int            // dimension of the embeddings
}

// NewEmbeddingIndex returns a new EmbeddingIndex.
func NewEmbeddingIndex() *EmbeddingIndex {
	return &EmbeddingIndex{
		Vocab: make(map[string]int),
	}
}

// LoadEmbeddings loads the token embeddings from a vocabulary file and a binary file containing the embeddings.
// The vocabulary file should be in JSON format and contain a map of tokens to their indices.
// The binary file should contain the following structure:
//
//	<num words>
//	<dimension>
//	<embeddings>
func (index *EmbeddingIndex) LoadEmbeddings(vocabPath, binPath string) error {
	defer utils.Measure("EmbeddingIndex_LoadEmbeddings")()

	vocabData, err := os.ReadFile(vocabPath)
	if err != nil {
		return err
	}

	// Parse the vocabulary.
	if err := json.Unmarshal(vocabData, &index.Vocab); err != nil {
		return err
	}

	// Load the embeddings.
	binFile, err := os.Open(binPath)
	if err != nil {
		return err
	}
	defer binFile.Close()

	var numWords, dim uint32

	err = binary.Read(binFile, binary.LittleEndian, &numWords)
	if err != nil {
		return err
	}

	err = binary.Read(binFile, binary.LittleEndian, &dim)
	if err != nil {
		return err
	}

	index.TokenEmbedings = make([]float64, int(numWords*dim))
	index.Dimension = int(dim)

	err = binary.Read(binFile, binary.LittleEndian, index.TokenEmbedings)
	if err != nil {
		return err
	}

	return nil
}

// Vector retrieves the embedding vector for the given word from the index.
// Returns all zero vectors for out-of-vocabulary words.
func (index *EmbeddingIndex) Vector(word string) []float64 {
	idx, ok := index.Vocab[word]
	if !ok {
		return make([]float64, index.Dimension)
	}

	start := idx * index.Dimension
	end := start + index.Dimension

	return index.TokenEmbedings[start:end]
}

// BuildFromDocuments builds and stores the embeddings of the given documents.
func (index *EmbeddingIndex) BuildFromDocuments(documents []string) error {
	defer utils.Measure("EmbeddingIndex_BuildFromDocuments")()

	index.DocumentEmbeddings = make([]float64, len(documents))
	for _, doc := range documents {
		embedding, err := index.EmbedDocument(doc)
		if err != nil {
			return err
		}
		index.DocumentEmbeddings = append(index.DocumentEmbeddings, embedding...)
	}

	return nil
}

// EmbedDocument embeds the given document and returns the embedding vector.
// Calculates the embedding by splitting the document into tokens and summing them.
// Out-of-vocabulary tokens are treated as all zero vectors.
func (index *EmbeddingIndex) EmbedDocument(document string) ([]float64, error) {
	if len(index.TokenEmbedings) == 0 {
		return nil, fmt.Errorf("empty token embeddings")
	}

	documentEmbedding := make([]float64, index.Dimension)

	tokens := utils.Tokenize(document)
	for _, token := range tokens {
		embedding := index.Vector(token)
		// Sum the token embeddings to get the document embedding.
		for i := 0; i < index.Dimension; i++ {
			documentEmbedding[i] += embedding[i]
		}
	}

	return documentEmbedding, nil
}

// CosineSimilarity calculates the cosine similarity between a vector `v` of shape [D]
// and all rows of a matrix `m` with shape [N x D].
// Returns a vector of shape [N] containing the cosine similarities.
func (index *EmbeddingIndex) CosineSimilarity(v, m []float64) ([]float64, error) {
	if len(v) != index.Dimension {
		return nil, fmt.Errorf("vector must have dimension %d", index.Dimension)
	}

	if len(m)%index.Dimension != 0 {
		return nil, fmt.Errorf("matrix must have a multiple of %d rows", index.Dimension)
	}

	// `n` is the number of rows in `m`, each with 'index.Dimension' columns.
	n := len(m) / index.Dimension
	result := make([]float64, n)

	// Normalize `v` and each row of `m`.
	vNorm := norm(v)
	for i := 0; i < n; i++ {
		start := i * index.Dimension
		end := start + index.Dimension
		row := m[start:end]

		rowNorm := norm(row)
		if rowNorm == 0 { // Cannot divide by zero.
			continue
		}

		// Calculate the dot product of `v` and `row`.
		var dot float64
		for j := 0; j < index.Dimension; j++ {
			dot += v[j] * row[j]
		}

		result[i] = dot / (vNorm * rowNorm)
	}

	return result, nil
}

// TopKNeighbors returns the indices of the k most similar documents to the given document.
func (index *EmbeddingIndex) TopKNeighbors(document string, k int) ([]int, error) {
	if len(index.DocumentEmbeddings) == 0 {
		return nil, fmt.Errorf("empty document embeddings")
	}

	documentEmbedding, err := index.EmbedDocument(document)
	if err != nil {
		return nil, err
	}

	cosineSimilarities, err := index.CosineSimilarity(documentEmbedding, index.DocumentEmbeddings)
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
func norm(v []float64) float64 {
	var norm float64
	for _, x := range v {
		norm += x * x
	}
	if norm == 0 {
		return 0
	}
	return math.Sqrt(norm)
}
