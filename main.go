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
	fmt.Println(cleanInput, len(cleanInput))

	snowballStemmer := func(input string, language string) (string, error) {
		return snowball.Stem(input, language, false)
	}

	var setOfLanguages = []string{"english", "russian"}
	result := stemInMultipleLanguages(snowballStemmer, cleanInput, setOfLanguages)
	fmt.Println(result)
}
