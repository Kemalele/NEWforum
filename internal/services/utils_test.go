package services

import (
	"testing"
)

func TestValidURL(t *testing.T) {
	assert := Assert(t, true)
	tests := []struct {
		url      string
		expected bool
		size     int
	}{
		{"/post/a4cfb14b-1c1e-4800-8afc-d6dadb6effe6", true, 2},
		{"/post/a4cfb14b-1c1e-4800-8afc-d6dadb6effe6/asdfg", false, 2},
		{"/write", true, 1},
	}

	for _, test := range tests {
		result := ValidURL(test.url, test.size)
		assert(result == test.expected, "expected:", test.expected, "got: ", result)
	}
}
