package words

import (
	sw "github.com/toadharvard/stopwords-iso"
)

// clearInputFromStopWords based on list of words from
// https://github.com/toadharvard/stopwords-iso/
//
// List of words probably will be changed
func clearInputFromStopWords(inputString []string) []string {
	var clearedStrings []string
	stopwordsMapping, _ := sw.NewStopwordsMapping()
	for _, str := range inputString {
		newStr := stopwordsMapping.ClearString(str)
		if newStr != "" {
			clearedStrings = append(clearedStrings, newStr)
		}
	}
	return clearedStrings
}
