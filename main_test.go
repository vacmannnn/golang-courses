package main

import (
	"fmt"
	kljSnowBall "github.com/kljensen/snowball"
	tbkSnowBall "github.com/tebeka/snowball"
	"strings"
	"testing"
)

const testString = "how to effectively bench golang stammers without any knowledge and experience ? " +
	"если бы я знал ответ..."

var splittedTestString = strings.Split(testString, " ")

func BenchmarkSnowballByKljensen(b *testing.B) {
	snowballStemmer := func(input string, language string) (string, error) {
		return kljSnowBall.Stem(input, language, false)
	}
	for i := 0; i < b.N; i++ {
		stemInMultipleLanguages(snowballStemmer, splittedTestString, []string{"english", "russian"})
	}
}

func BenchmarkSnowballByTebeka(b *testing.B) {
	snowballStemmer := func(input string, language string) (string, error) {
		myStemmer, err := tbkSnowBall.New(language)
		if err != nil {
			fmt.Println("error", err)
			return "", nil
		}
		defer myStemmer.Close()
		return myStemmer.Stem(input), nil
	}
	for i := 0; i < b.N; i++ {
		stemInMultipleLanguages(snowballStemmer, splittedTestString, []string{"english", "russian"})
	}
}
