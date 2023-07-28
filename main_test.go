package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToString(t *testing.T) {

	tests := []struct {
		input  interface{}
		expect string
	}{
		{
			input:  "RDS-EVENT-0042",
			expect: "RDS-EVENT-0042",
		},
		{
			input:  "",
			expect: "",
		},
		{
			input:  133,
			expect: "133",
		},
		{
			input:  true,
			expect: "true",
		},
	}

	for _, test := range tests {
		output := ConvertToString(test.input)
		assert.IsType(t, test.expect, output)
		assert.Equal(t, test.expect, output)
	}
}

func TestGetUniqueSnapshot(t *testing.T) {

	tests := []struct {
		target []string
		source []string
		expect []string
		err    error
	}{
		{
			target: []string{"a", "b", "c"},
			source: []string{"a"},
			expect: []string(nil),
			err:    nil,
		},
		{
			target: []string{"a", "b", "c"},
			source: []string{"a", "b", "c"},
			expect: []string(nil),
			err:    nil,
		},
		{
			target: []string{"a"},
			source: []string{"a", "b", "c"},
			expect: []string{"b", "c"},
			err:    nil,
		},
	}

	for _, test := range tests {
		output, _ := GetUniqueSnapShots(test.target, test.source)
		assert.IsType(t, []string{}, output)
		assert.Equal(t, test.expect, output)
	}
}
