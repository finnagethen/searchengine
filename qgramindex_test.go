package main

import (
	"reflect"
	"testing"
)

func TestBuildFormFile(t *testing.T) {
	index := NewQGramIndex(3)
	err := index.BuildFormFile("test.tsv")
	if err != nil {
		t.Fatalf("BuildFormFile failed: %v", err)
	}

	expectedInfos := []Info{
		{
			Name:  "frei",
			Score: 3,
			Infos: []string{"first entity", "used for doctests"},
		},
		{
			Name:  "brei",
			Score: 2,
			Infos: []string{"second entity", "also for doctests"},
		},
	}

	// Synonyms: "frei"(0), "fre"(0), "fri"(0), "brei"(1), "bre"(1), "bri"(1)
	expectedSynonymToRecord := []int{0, 0, 0, 1, 1, 1}

	expectedInvertedLists := map[string][]Posting{
		"$$f": {
			{ID: 0, Frequency: 1}, {ID: 1, Frequency: 1}, {ID: 2, Frequency: 1},
		},
		"$fr": {
			{ID: 0, Frequency: 1}, {ID: 1, Frequency: 1}, {ID: 2, Frequency: 1},
		},
		"fre": {
			{ID: 0, Frequency: 1}, {ID: 1, Frequency: 1},
		},
		"rei": {
			{ID: 0, Frequency: 1}, {ID: 3, Frequency: 1},
		},
		"fri": {
			{ID: 2, Frequency: 1},
		},
		"$$b": {
			{ID: 3, Frequency: 1}, {ID: 4, Frequency: 1}, {ID: 5, Frequency: 1},
		},
		"$br": {
			{ID: 3, Frequency: 1}, {ID: 4, Frequency: 1}, {ID: 5, Frequency: 1},
		},
		"bre": {
			{ID: 3, Frequency: 1}, {ID: 4, Frequency: 1},
		},
		"bri": {
			{ID: 5, Frequency: 1},
		},
	}

	if !reflect.DeepEqual(index.Infos, expectedInfos) {
		t.Errorf("Infos mismatch.\nExpected: %+v\nGot: %+v",
			expectedInfos, index.Infos)
	}

	if !reflect.DeepEqual(index.SynonymToRecord, expectedSynonymToRecord) {
		t.Errorf("SynonymToRecord mismatch.\nExpected: %+v\nGot: %+v",
			expectedSynonymToRecord, index.SynonymToRecord)
	}

	if !reflect.DeepEqual(index.InvertedLists, expectedInvertedLists) {
		t.Errorf("InvertedLists mismatch.\nExpected: %+v\nGot: %+v",
			expectedInvertedLists, index.InvertedLists)
	}
}
