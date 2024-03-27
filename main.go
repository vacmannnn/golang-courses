package main

import (
	"fmt"
	"github.com/kljensen/snowball"
)

func main() {
	inputString, err := getStringFromArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	var cleanInput = clearInputFromStopWords(inputString)
	fmt.Printf("cleared input -- %s\n", cleanInput)

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	var setOfLanguages = []string{"english", "russian"}
	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)
	fmt.Printf("stemmed and cleared input -- %s\n", result)
}
