package qgramindex

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/finnagethen/searchengine/internal/ped"
	"github.com/finnagethen/searchengine/internal/utils"
)

// RecordID defines an ID type for records
type RecordID int

// SynonymID defines an ID type for synonyms
type SynonymID int

// Posting defines an entry in an inverted list
type Posting struct {
	ID        SynonymID
	Frequency int
}

// Info defines the information for each record (document)
type Info struct {
	Name  string
	Score int
	Infos []string
}

// Match defines the return type of prefix querys
type Match struct {
	ID  RecordID
	PED int
}

type QGramIndex struct {
	Q               int
	InvertedLists   map[string][]Posting // q-gram -> posting list
	SynonymToRecord []RecordID           // synonym id -> record id
	NormedNames     []string             // synonym id -> normalized synonym
	Infos           []Info               // record id -> info
}

// NewQGramIndex returns a new QGramIndex.
func NewQGramIndex(q int) *QGramIndex {
	return &QGramIndex{
		Q:             q,
		InvertedLists: make(map[string][]Posting),
	}
}

// BuildFormFile builds an index from a file.
// The file is expected to contain one record per line, in the format:
//
//	<name>\t<score>\t<synonyms>\t<info1>\t<info2>\t...
//
// Semicolons separate synonyms. An example record:
//
//	Albert Einstein\t275\tEinstein;A. Einstein\tGerman physicist\t...
func (index *QGramIndex) BuildFormFile(path string) error {
	defer utils.Measure("QGramIndex_BuildFormFile")()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var recordID RecordID
	var synonymID SynonymID

	scanner := bufio.NewScanner(file)
	scanner.Scan() // skip the header

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		name, score, synonyms, infos := fields[0], fields[1], fields[2], fields[3:]

		scoreConverted, err := strconv.Atoi(score)
		if err != nil {
			return err
		}

		// Cache the name, score and additional info.
		index.Infos = append(index.Infos, Info{
			Name:  name,
			Score: scoreConverted,
			Infos: infos,
		})

		// Calculate the q-grams for every name.
		names := append([]string{name}, strings.Split(synonyms, ";")...)
		for _, n := range names {
			index.SynonymToRecord = append(index.SynonymToRecord, recordID)
			normedName := utils.Normalize(n)
			index.NormedNames = append(index.NormedNames, normedName)

			for _, qgram := range computeQGrams(normedName, index.Q) {
				postings := index.InvertedLists[qgram]
				postingsLen := len(postings)

				if postingsLen > 0 && postings[postingsLen-1].ID == synonymID {
					postings[postingsLen-1].Frequency++
				} else {
					index.InvertedLists[qgram] = append(postings, Posting{
						ID:        synonymID,
						Frequency: 1,
					})
				}
			}
			synonymID++
		}
		recordID++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// FindMatches retrieves all postings with PED(x, y) <= delta for a given integer delta
// and prefix x. The prefix should be normalized and non-empty.
// Returns a list of (ID, PED) tuples ordered first by PED and then record score.
func (index *QGramIndex) FindMatches(prefix string, delta int) ([]Match, error) {
	defer utils.Measure("QGramIndex_FindMatches")()

	if len(prefix) == 0 {
		return nil, fmt.Errorf("prefix must not be empty")
	}

	threshold := len(prefix) - (index.Q * delta)

	if threshold <= 0 {
		return nil, fmt.Errorf("threshold must be positive (got %d); adjust delta", threshold)
	}

	// Count frequencies of qgrams in the prefix.
	prefixQGrams := computeQGrams(prefix, index.Q)
	prefixFreqs := make(map[string]int)
	for _, q := range prefixQGrams {
		prefixFreqs[q]++
	}

	// Maps the synonyms to the number of qgrams in common with the prefix.
	candidates := make(map[SynonymID]int)
	for qgram, freq := range prefixFreqs {
		postings := index.InvertedLists[qgram]
		for _, posting := range postings {
			candidates[posting.ID] += min(posting.Frequency, freq)
		}
	}

	// Filter out synonyms with less than `threshold` qgrams in common.
	// Calculate the PED(x, y) for all relevant synonyms.
	matches := make(map[RecordID]int)
	pedCalculations := 0 // Counter for stats.
	for synID, freq := range candidates {
		if freq < threshold {
			continue
		}

		distance := ped.PrefixEditDistance(prefix, index.NormedNames[synID], delta)
		pedCalculations++
		if distance <= delta {
			recID := index.SynonymToRecord[synID]

			// If a record has multiple matching synonyms, keep the one with the lowest PED.
			if oldDistance, ok := matches[recID]; !ok || distance < oldDistance {
				matches[recID] = distance
			}
		}
	}

	slog.Info("ped calculations",
		"completed", pedCalculations,
		"candidates", len(candidates),
	)
	// Convert the map to a slice of (ID, PED) tuples.
	var result []Match
	for recID, bestPED := range matches {
		result = append(result, Match{ID: recID, PED: bestPED})
	}

	// Sort first by PED (ascending) and then by record-score (descending).
	slices.SortFunc(result, func(a, b Match) int {
		pedCmp := a.PED - b.PED // smaller is better
		if pedCmp != 0 {
			return pedCmp
		}
		return index.Infos[b.ID].Score - index.Infos[a.ID].Score // bigger is better
	})

	return result, nil
}

// GetInfo retrieves the info for a given synonym ID.
func (index *QGramIndex) GetInfo(id SynonymID) Info {
	recordID := index.SynonymToRecord[id]
	return index.Infos[recordID]
}

// computeQGrams computes the q-grams for the given word.
func computeQGrams(word string, q int) []string {
	padding := strings.Repeat("$", q-1)
	padded := padding + word

	qgrams := make([]string, 0, len(word))
	for i := 0; i < len(word); i++ {
		qgrams = append(qgrams, padded[i:i+q])
	}

	return qgrams
}
