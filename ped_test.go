package main

import "testing"

func TestPrefixEditDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"both_empty", "", "", 0},
		{"first_empty", "", "abc", 0},
		{"second_empty", "abc", "", 3},
		{"identical", "abc", "abc", 0},
		{"s2_prefix_of_s1", "prefix", "pre", 3},
		{"s1_prefix_of_s2", "pre", "prefix", 0},
		{"close_prefix", "cwit", "wltbubhcp", 2},
		{"different", "abc", "xyz", 3},
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
