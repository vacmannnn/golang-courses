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

	// TODO: может ли быть больше 1 аргумента ? проверить случай с \n
	return strings.Split(inputString[0], " "), nil
}

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
