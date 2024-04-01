package words

import (
	"fmt"
	"github.com/kljensen/snowball"
	"strings"
)

// StemStringWithClearing
// TODO: comments on public func
func StemStringWithClearing(input string) []string {

	splittedString := strings.Split(input, " ")

	var cleanInput = clearInputFromStopWords(splittedString)
	fmt.Printf("cleared input -- %s\n", cleanInput)

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	var setOfLanguages = []string{"english", "russian"}
	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)
	fmt.Printf("stemmed and cleared input -- %s\n", result)

	return result
}
