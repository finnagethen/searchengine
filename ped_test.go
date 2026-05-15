package main

import "testing"

func TestEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"Both strings empty", "", "", 0},
		{"First string empty", "", "abc", 3},
		{"Second string empty", "abc", "", 3},
		{"Identical strings", "abc", "abc", 0},

		{"Single substitution", "abc", "axc", 1},
		{"Single insertion", "abc", "abxc", 1},
		{"Single deletion", "abc", "ac", 1},
		{"Completely different", "abc", "xyz", 3},

		{"Kitten/Sitting", "kitten", "sitting", 3},
		{"Flaw/Lawn", "flaw", "lawn", 2},
		{"Different lengths", "wltbubhcp", "cwit", 8},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := EditDistance(tc.s1, tc.s2)
			if got != tc.expected {
				t.Errorf("EditDistance(%q, %q) = %d; want %d", tc.s1, tc.s2, got, tc.expected)
			}
		})
	}
}

func TestPrefixEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"Both strings empty", "", "", 0},
		{"First string empty", "", "abc", 0},
		{"Second string empty", "abc", "", 3},
		{"Identical strings", "abc", "abc", 0},
		{"s2 is prefix of s1", "prefix", "pre", 3},
		{"s1 is prefix of s2", "pre", "prefix", 0},
		{"Partial match close prefix", "cwit", "wltbubhcp", 2},
		{"Completely different", "abc", "xyz", 3},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PrefixEditDistance(tc.s1, tc.s2)
			if got != tc.expected {
				t.Errorf("PrefixEditDistance(%q, %q) = %d; want %d", tc.s1, tc.s2, got, tc.expected)
			}
		})
	}
}
