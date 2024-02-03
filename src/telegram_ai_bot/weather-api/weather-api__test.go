package weather_api

import "testing"

type FitToMaxSpaceTestingArgs struct {
	expected, initial string
	maxChars          int
}

func TestFitToMaxSpace(t *testing.T) {

	args := []FitToMaxSpaceTestingArgs{
		{expected: " 10:34 ", initial: "10:34", maxChars: 7},
		{expected: "    -1     ", initial: "-1", maxChars: 11},
	}

	for _, arg := range args {
		actual := fitToMaxSpace(arg.initial, arg.maxChars)
		if actual != arg.expected {
			t.Logf(`Expected '%s', but got '%s'`, arg.expected, actual)
			t.Fail()
		}
	}
}
