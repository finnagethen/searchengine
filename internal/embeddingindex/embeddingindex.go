package embeddingindex

type EmbeddingIndex struct {
	Vocab              map[string]int // token -> token index
	TokenEmbedings     []float32      // token index -> embedding
	DocumentEmbeddings []float32      // document index -> embedding
}

func NewEmbeddingIndex() *EmbeddingIndex {
	return &EmbeddingIndex{
		Vocab: make(map[string]int),
	}
}

func (index *EmbeddingIndex) LoadEmbeddings(path string) error {
	return nil
}

func (index *EmbeddingIndex) BuildFromDocuments(documents []string) error {
	return nil
}

func (index *EmbeddingIndex) EmbedDocument(document string) []float32 {
	return nil
}

func (index *EmbeddingIndex) CosineSimilarity(a, b []float32) float32 {
	return 0
}

func (index *EmbeddingIndex) TopKNeighbors(document []float32, k int) []int {
	return nil
}
