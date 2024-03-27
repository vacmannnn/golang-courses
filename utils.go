package main

import (
	"errors"
	"flag"
	sw "github.com/toadharvard/stopwords-iso"
	"strings"
)

func getStringFromArguments() ([]string, error) {
	useString := flag.Bool("s", false, "get input string")
	flag.Parse()

	if !(*useString) {
		return []string{}, errors.New("expected for -s flag")
	}

	inputString := flag.Args()
	if len(inputString) == 0 {
		return []string{}, errors.New("expected for non-empty string")
	}

	// avoid case with other flags and strings
	if len(inputString) > 1 {
		return []string{}, errors.New("expected format is \"./myapp -s string_to_stem\"")
	}

	return strings.Split(inputString[0], " "), nil
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
