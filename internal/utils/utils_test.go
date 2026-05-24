package utils

import (
	"reflect"
	"testing"
)

func TestIsAlphanumeric(t *testing.T) {
	tests := []struct {
		char     byte
		expected bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'0', true},
		{'9', true},

		{'!', false},
		{'-', false},
		{'"', false},
		{' ', false},
		{'\n', false},
	}

	for _, tc := range tests {
		got := IsAlphanumeric(tc.char)

		if got != tc.expected {
			t.Errorf("IsAlphanumeric(%q) = %v; want %v",
				tc.char, got, tc.expected)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"HELLO", "hello"},
		{"Hello123", "hello123"},
		{"Hello, World!", "helloworld"},
		{"don't", "dont"},
		{"mother-in-law", "motherinlaw"},
		{"   test   ", "test"},
		{"", ""},
		{"!!!", ""},
	}

	for _, tc := range tests {
		got := Normalize(tc.input)

		if got != tc.expected {
			t.Errorf("Normalize(%q) = %q; want %q",
				tc.input, got, tc.expected)
		}
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			"Hello world",
			[]string{"Hello", "world"},
		},
		{
			"don't stop-believing",
			[]string{"don't", "stop-believing"},
		},
		{
			"Hello, world!",
			[]string{"Hello", "world"},
		},
		{
			"rock'n'roll",
			[]string{"rock'n'roll"},
		},
		{
			"",
			nil,
		},
		{
			"!!!",
			nil,
		},
	}

	for _, tc := range tests {
		got := Tokenize(tc.input)

		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("Tokanize(%q) = %#v; want %#v",
				tc.input, got, tc.expected)
		}
	}
}

func TestEqualSliceEpsilon(t *testing.T) {
	tests := []struct {
		a, b     []float32
		epsilon  float32
		expected bool
	}{
		{[]float32{1.0, 2.0}, []float32{1.0, 2.0}, 0.001, true},
		{[]float32{1.0, 2.0}, []float32{1.001, 2.001}, 0.01, true},
		{[]float32{1.0, 2.0}, []float32{1.01, 2.01}, 0.001, false},
		{[]float32{1.0}, []float32{1.0, 2.0}, 0.001, false},
	}

	for _, tc := range tests {
		got := EqualSliceEpsilon(tc.a, tc.b, tc.epsilon)

		if got != tc.expected {
			t.Errorf("EqualSliceEpsilon(%v, %v, %f) = %v; want %v",
				tc.a, tc.b, tc.epsilon, got, tc.expected)
		}
	}
}
