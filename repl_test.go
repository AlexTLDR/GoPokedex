package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "HELLO WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  foo   bar   baz   ",
			expected: []string{"foo", "bar", "baz"},
		},
		{
			input:    "one",
			expected: []string{"one"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "123 456 789",
			expected: []string{"123", "456", "789"},
		},
		{
			input:    "  mixed   CASE   input  ",
			expected: []string{"mixed", "case", "input"},
		},
		{
			input:    "  leading   spaces  ",
			expected: []string{"leading", "spaces"},
		},
		{
			input:    "trailing   spaces   ",
			expected: []string{"trailing", "spaces"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("lengths don't match: %v vs %v", actual, c.expected)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q) = %q, expected %q", c.input, word, expectedWord)
			}
		}
	}
}
