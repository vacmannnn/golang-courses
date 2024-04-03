package words

import (
	"github.com/kljensen/snowball"
	"strings"
)

// StemStringWithClearing gets string, clears from most popular words and returns slice of keys
func StemStringWithClearing(input string) []string {

	splittedString := strings.Split(input, " ")

	var cleanInput = clearInputFromStopWords(splittedString)

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	var setOfLanguages = []string{"english", "russian"}
	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)

	return result
}
