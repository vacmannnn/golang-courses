package words

import (
	sw "github.com/toadharvard/stopwords-iso"
	"regexp"
	"strings"
)

// clearInput removes stopwords and marks as "'", "?", "-".
// stopwords cleaner based on https://github.com/toadharvard/stopwords-iso/
//
// List of words probably will be changed
func clearInput(inputString []string) []string {
	joinedString := strings.Join(inputString, " ")
	tokenizedString := tokenize(joinedString)

	stopWordsMapping, _ := sw.NewStopwordsMapping()
	var clearedString []string

	for _, str := range tokenizedString {
		newStr := stopWordsMapping.ClearString(str)
		if len(newStr) > 2 {
			clearedString = append(clearedString, newStr)
		}
	}
	return clearedString
}

func tokenize(str string) []string {
	wordSegmenter := regexp.MustCompile(`[\pL\p{Mc}\p{Mn}-_']+`)
	alphanumericOnly := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	words := strings.Fields(str)
	splitted := strings.Join(words, " ")

	onlyAlphanumeric := alphanumericOnly.ReplaceAllString(splitted, " ")
	words = wordSegmenter.FindAllString(onlyAlphanumeric, -1)
	return words
}
