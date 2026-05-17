// Function to calculate the prefix edit distance of two strings.
package main

import "slices"

// PrefixEditDistance returns the PED of two strings, where `s1` is a prefix of `s2`
// if the PED is smaller or equal to `delta`; `delta` + 1 otherwise.
// Assume `delta` >= 0.
func PrefixEditDistance(s1, s2 string, delta int) int {
	// defer MeasureExecutionTime("PrefixEditDistance")()

	n := len(s1) + 1
	// It's enough to compute the first |s1| + delta + 1 columns.
	m := min(n+delta, len(s2)+1)

	// If `s2`is shorter than `s1` - delta, then the PED cannot be smaller or equal to `delta`.
	if m < n-delta {
		return delta + 1
	}

	// Only store the current and previous row and not the whole matrix.
	// Initialize the first row to 0...m-1.
	previousRow := make([]int, m)
	currentRow := make([]int, m)
	for col := 0; col < m; col++ {
		previousRow[col] = col
	}

	// Calculate the edit distance by taking the minimum cost of the three possible operations.
	for row := 1; row < n; row++ {
		currentRow[0] = row
		for col := 1; col < m; col++ {
			cost := func() int {
				isDiffrent := 1
				if s1[row-1] == s2[col-1] {
					isDiffrent = 0
				}
				return min(previousRow[col-1]+isDiffrent, previousRow[col]+1, currentRow[col-1]+1)
			}
			currentRow[col] = cost()
		}

		previousRow, currentRow = currentRow, previousRow
	}

	// Min of `previousRow` since it was swapt with `currentRow` in the loop.
	return min(slices.Min(previousRow), delta+1)
}
