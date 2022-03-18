package learn_testing

import (
	"testing"
)

type TestCase struct {
	value    int
	expected bool
	actual   bool
}

func TestCalculateArmstrongNumber(t *testing.T) {
	t.Run("should return true for value of 371", func(t *testing.T) {
		testCase := TestCase{
			value:    371,
			expected: true,
		}
		testCase.actual = isArmstrongNumber(testCase.value)
		if testCase.actual != testCase.expected {
			t.Fail()
		}
	})
	t.Run("should return true for value of 370", func(t *testing.T) {
		testCase := TestCase{
			value:    370,
			expected: true,
		}
		testCase.actual = isArmstrongNumber(testCase.value)
		if testCase.actual != testCase.expected {
			t.Fail()
		}
	})
}
func TestNegativeCalculateArmstrongNumber(t *testing.T) {
	testCase := TestCase{
		value:    350,
		expected: false,
	}
	testCase.actual = isArmstrongNumber(testCase.value)
	if testCase.actual != testCase.expected {
		t.Fail()
	}
}
