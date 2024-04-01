package words

import (
	"errors"
	"flag"
	sw "github.com/toadharvard/stopwords-iso"
	"strings"
)

func getStringFromArguments() ([]string, error) {
	var inputString string
	flag.StringVar(&inputString, "s", "", "get input string")
	flag.Parse()

	otherInput := flag.Args()

	// avoid case with other flags and strings
	if len(otherInput) > 0 || inputString == "" {
		return []string{}, errors.New("expected format is \"./myapp -s stringToStem\"")
	}
	return strings.Split(inputString, " "), nil
}

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
