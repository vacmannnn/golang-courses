package main

import (
	"fmt"
	"github.com/kljensen/snowball"
)

var setOfLanguages = []string{"english", "russian"}

func main() {
	inputString, err := getStringFromArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	var cleanInput = clearInputFromStopWords(inputString)
	fmt.Println(cleanInput, len(cleanInput))

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)
	fmt.Println(result)
}
