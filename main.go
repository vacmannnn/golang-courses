package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	sw "github.com/toadharvard/stopwords-iso"
	"strings"
)

type Stemmer func(input string, language string) (string, error)

var setOfLanguages = []string{"english", "russian"}

func main() {
	inputString, err := getStringFromArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	var cleanInput = clearInput(inputString)
	fmt.Println(cleanInput, len(cleanInput))

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)
	fmt.Println(result)
}

func stemInMultipleLanguages(myStemmer Stemmer, input []string, languages []string) []string {
	var output []string
	for _, language := range languages {
		output = stemSingleLanguage(myStemmer, input, language)
		input = output
	}
	return output
}

func stemSingleLanguage(myStemmer Stemmer, input []string, language string) []string {
	res := make([]string, 0, len(input))
	for _, str := range input {
		newStr, err := myStemmer(str, language)
		if err != nil {
		} else {
			res = append(res, newStr)
		}
	}
	return removeDuplicateStrings(res)
}

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

func clearInput(inputString []string) []string {
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

func removeDuplicateStrings(strings []string) []string {
	var list []string
	keys := make(map[string]bool)

	for _, str := range strings {
		if _, value := keys[str]; !value {
			keys[str] = true
			list = append(list, str)
		}
	}
	return list
}
