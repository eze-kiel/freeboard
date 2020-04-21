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

func TestAuthorizedURL(t *testing.T) {
	ts := []struct {
		input    string
		expected bool
	}{
		{"www.google.com", true},
		{"facebook.com", true},
		{"http://www.bit.ly", false},
	}

	for _, tc := range ts {
		t.Run(tc.input, func(t *testing.T) {
			res := AuthorizedURL(tc.input)
			if res != tc.expected {
				t.Errorf("got %v, expected %v", res, tc.expected)
			}
		})
	}
}

func TestAuthorizedText(t *testing.T) {
	ts := []struct {
		input    string
		expected bool
	}{
		{"hey boi", true},
		{"how are you", true},
		{"sentence to run tests", false},
	}

	for _, tc := range ts {
		t.Run(tc.input, func(t *testing.T) {
			res := AuthorizedText(tc.input)
			if res != tc.expected {
				t.Errorf("got %v, expected %v", res, tc.expected)
			}
		})
	}
}

func TestCheckCategory(t *testing.T) {
	ts := []struct {
		input    string
		expected bool
	}{
		{"nature", true},
		{"all", true},
		{"science", true},
		{"", false},
		{"alchemy", false},
	}

	for _, tc := range ts {
		t.Run(tc.input, func(t *testing.T) {
			res := CheckCategory(tc.input)
			if res != tc.expected {
				t.Errorf("got %v, expected %v", res, tc.expected)
			}
		})
	}
}
