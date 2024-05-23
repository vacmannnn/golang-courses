package words

import (
	"reflect"
	"strings"
	"testing"
)

var cases1 = []struct {
	name     string
	input    string
	expected []string
}{
	{
		name:     "hello world",
		input:    "hello, world !",
		expected: []string{},
	},
	{
		name:     "simple marks",
		input:    "am i, the? greatest! tester?!?!?!",
		expected: strings.Split("tester", " "),
	},
	{
		name:     "different marks with whitespaces",
		input:    "impossible         !!test?!@@ 		case with-my O'connor--- thing and $$$",
		expected: strings.Split("imposs connor", " "),
	},
	{
		name:     "Apple doctor",
		input:    "apple,doctor",
		expected: strings.Split("appl doctor", " "),
	},
}

func TestStemmer(t *testing.T) {
	for _, tc := range cases1 {
		t.Run(tc.name, func(t *testing.T) {
			res := StemStringWithClearing(strings.Split(tc.input, " "))
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("got %v, want %v", res, tc.expected)
			}
		})
	}
}
