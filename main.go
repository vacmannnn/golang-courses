package main

import (
	"flag"
	"fmt"
	"github.com/bbalet/stopwords"
)

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

	var clearedStrings = clearInput(inputString)
	fmt.Println(clearedStrings, len(clearedStrings))
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
