package utils_test

import (
	"testing"

	"github.com/berberapan/dota-work/internal/utils"
)

func TestConvertIntToBinary(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{
			input:    0,
			expected: "0",
		},
		{
			input:    63,
			expected: "111111",
		},
		{
			input:    1542,
			expected: "11000000110",
		},
	}

	for _, test := range testCases {
		actual := utils.ConvertIntToBinary(test.input)
		if actual != test.expected {
			t.Errorf("\nReturn values not matching.\nActual: %s Expected: %s", actual, test.expected)
		}
	}
}

func TestConvertDateToUnix(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{
			input:    "2019-11-02",
			expected: 1572652800,
		},
		{
			input:    "2024-02-29",
			expected: 1709164800,
		},
		{
			input:    "2050-01-01",
			expected: 2524608000,
		},
	}

	for _, test := range testCases {
		actual := utils.ConvertDateToUnix(test.input)
		if actual != test.expected {
			t.Errorf("\nReturn values not matching.\nActual: %d Expected: %d", actual, test.expected)
		}
	}
}
