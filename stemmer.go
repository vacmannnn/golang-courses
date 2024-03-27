package main

type stemmer func(input string, language string) (string, error)

func stemInMultipleLanguages(myStemmer stemmer, input []string, languages []string) []string {
	var output []string
	for _, language := range languages {
		output = stemSingleLanguage(myStemmer, input, language)
		input = output
	}
	return output
}

func stemSingleLanguage(myStemmer stemmer, input []string, language string) []string {
	res := make([]string, 0, len(input))
	for _, str := range input {
		newStr, err := myStemmer(str, language)
		// TODO: error handling
		if err != nil {
		} else {
			res = append(res, newStr)
		}
	}
	return removeDuplicateStrings(res)
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
