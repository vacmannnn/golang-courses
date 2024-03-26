package main

import (
	"flag"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/kljensen/snowball"
	"strings"
)

type Stemmer func(input string, language string) (string, error)

func main() {
	useString := flag.Bool("s", false, "get input string")
	flag.Parse()

	if !(*useString) {
		fmt.Println("expected for -s flag")
		return
	}

	inputString := flag.Args()
	if len(inputString) == 0 {
		fmt.Println("expected for non-empty string")
		return
	}

	inputString = strings.Split(inputString[0], " ")
	var clearedStrings = clearInput(inputString)
	fmt.Println(clearedStrings, len(clearedStrings))

	myStem := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	result := stemInput(myStem, clearedStrings)
	fmt.Println(result)
}

func stemInput(myStemmer Stemmer, input []string) []string {
	res := make([]string, 0, len(input))
	for _, str := range input {
		newStr, err := myStemmer(str, "english")
		if err != nil {

		} else {
			res = append(res, newStr)
		}
	}
	return removeDuplicateStrings(res)
}

func clearInput(inputString []string) []string {
	var clearedStrings []string
	for _, str := range inputString {
		newStr := stopwords.CleanString(str, "en", false)
		if newStr != " " {
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
