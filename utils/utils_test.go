package utils

import "testing"

func TestIsURL(t *testing.T) {
	ts := []struct {
		input    string
		expected bool
	}{
		{"blabla", false},
		{"", false},
		{"www.google.com", false},
		{"http://www.google.com", true},
		{"http://stackoverflow.com/", true},
	}

	for _, tc := range ts {
		t.Run(tc.input, func(t *testing.T) {
			res := IsURL(tc.input)
			if res != tc.expected {
				t.Errorf("got %v, expected %v", res, tc.expected)
			}
		})
	}
}
