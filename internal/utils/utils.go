package utils

import (
	"log/slog"
	"math"
	"strings"
	"time"
)

// Measure returns a function that measures the execution time
// of the wrapped function.
func Measure(name string) func() {
	start := time.Now()
	return func() {
		slog.Info("finished",
			"name", name,
			"duration", time.Since(start),
		)
	}
}

// IsAlphanumeric returns true if the given character is an alphanumeric
// character.
func IsAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

// Normalize normalizes a string to lower case and removes all
// non-alphanumeric characters.
func Normalize(word string) string {
	var builder strings.Builder
	for i := 0; i < len(word); i++ {
		c := word[i]
		if IsAlphanumeric(c) {
			builder.WriteByte(c)
		}
	}

	return strings.ToLower(builder.String())
}

// EqualSliceEpsilon checks if two float32 slices are approximately equal.
func EqualSliceEpsilon(a, b []float32, epsilon float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if float32(math.Abs(float64(a[i]-b[i]))) > epsilon {
			return false
		}
	}
	return true
}
