// A simple qgram-index.
package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
	ID  SynonymID
	PED int
}

type QGramIndex struct {
	Q               int
	InvertedLists   map[string][]Posting // q-gram -> posting list
	SynonymToRecord []RecordID           // synonym id -> record id
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
			normedName := normalize(n)
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
	if len(prefix) == 0 {
		return nil, fmt.Errorf("prefix must not be empty")
	}

	threshold := len(prefix) - (index.Q * delta)

	if threshold <= 0 {
		return nil, fmt.Errorf("threshold must be positive (got %d); adjust delta", threshold)
	}

	// Maps the records to the number of qgrams in common with the prefix.
	candidates := make(map[SynonymID]int)
	for _, qgram := range computeQGrams(prefix, index.Q) {
		// Find all records with at least `threshold` qgrams in common.
		matches := index.InvertedLists[qgram]
		for _, match := range matches {
			candidates[match.ID] += 1
		}
	}

	// Filter out records with less than `threshold` qgrams in common.
	// Calculate the PED(x, y) for all relevant records.
	var result []Match
	for id, matchCount := range candidates {
		if matchCount < threshold {
			result = append(result, Match{id, matchCount})
		}
	}

	// Sort by PED and then record-score descending.
	slices.SortFunc(result, func(a, b Match) int {
		pedCmp := b.PED - a.PED
		if pedCmp != 0 {
			return pedCmp
		}
		return index.Infos[b.ID].Score - index.Infos[a.ID].Score
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

// isAlphanumeric returns true if the given character is an alphanumeric
// character.
func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// normalize normalizes a string to lower case and removes all
// non-alphanumeric characters.
func normalize(word string) string {
	var builder strings.Builder
	for i := 0; i < len(word); i++ {
		c := word[i]
		if isAlphanumeric(c) {
			builder.WriteByte(c)
		}
	}

	return strings.ToLower(builder.String())
}
