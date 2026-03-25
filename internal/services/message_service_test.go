package services

import "testing"

func TestNormalizeSenderName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "alexis lowercase", input: "alexis", expected: "Alexis"},
		{name: "alexis with suffix", input: "alexis :)", expected: "Alexis"},
		{name: "alexis mixed case", input: "ALeXiS", expected: "Alexis"},
		{name: "anyel lowercase", input: "anyel", expected: "Anyel"},
		{name: "anyel with spaces", input: "  anyel  ", expected: "Anyel"},
		{name: "anyel mixed case", input: "AnYeL", expected: "Anyel"},
		{name: "unknown user", input: "Carlos", expected: "Carlos"},
		{name: "unknown trimmed", input: "  Laura  ", expected: "Laura"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeSenderName(tt.input)
			if got != tt.expected {
				t.Fatalf("resultado inesperado: esperado %q, obtenido %q", tt.expected, got)
			}
		})
	}
}
