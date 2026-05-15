// Functions to calculate the edit and prefix edit distance of two strings.
package main

import "slices"

// editDistanceLastRow calculates the last row of the edit distance matrix using a space-optimized approach by storing only two rows.
func editDistanceLastRow(s1, s2 string) []int {
	n := len(s1)
	m := len(s2)

	// Only store the current and previous row and not the whole matrix.
	// Initialize the first row to 0...n.
	previousRow := make([]int, n+1)
	currentRow := make([]int, n+1)
	for col := 0; col <= n; col++ {
		previousRow[col] = col
	}

	// Calculate the edit distance by taking the minimum cost of the three possible operations.
	for row := 1; row <= m; row++ {
		currentRow[0] = row
		for col := 1; col <= n; col++ {
			cost := func() int {
				isDiffrent := 1
				if s1[col-1] == s2[row-1] {
					isDiffrent = 0
				}
				return min(previousRow[col-1]+isDiffrent, previousRow[col]+1, currentRow[col-1]+1)
			}
			currentRow[col] = cost()
		}

		previousRow, currentRow = currentRow, previousRow
	}

	// Return `previousRow` since it was swapt with `currentRow` in the loop.
	return previousRow
}

// EditDistance calculates the levenshtein distance between two strings.
func EditDistance(s1, s2 string) int {
	lastRow := editDistanceLastRow(s1, s2)
	return lastRow[len(s1)]
}

// PrefixEditDistance calculates the minimum edit distance between a prefix of s1 and s2.
func PrefixEditDistance(s1, s2 string) int {
	return slices.Min(editDistanceLastRow(s2, s1))
}
