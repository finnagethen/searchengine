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

func TestTokanize(t *testing.T) {
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
		got := Tokanize(tc.input)

		if !reflect.DeepEqual(got, tc.expected) {
			t.Errorf("Tokanize(%q) = %#v; want %#v",
				tc.input, got, tc.expected)
		}
	}
}
