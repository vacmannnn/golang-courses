package main

import (
	"cmp"
	"courses/internal/core"
	"slices"
)

func findByComics(comics map[int]core.ComicsDescript, input []string) []goodComics {
	var res []goodComics
	for id, v := range comics {
		var numOfWords int
		for _, word := range input {
			if slices.Contains(v.Keywords, word) {
				numOfWords++
			}
		}
		res = append(res, goodComics{id: id, numOfKeywords: numOfWords})
	}
	slices.SortFunc(res, func(a, b goodComics) int {
		return cmp.Compare(a.numOfKeywords, b.numOfKeywords) * (-1)
	})
	return res
}

func findByIndex(index map[string][]int, input []string) []goodComics {
	wasFound := make(map[int]int)
	for _, keywords := range input {
		for _, comicsID := range index[keywords] {
			wasFound[comicsID]++
		}
	}
	var res []goodComics
	for k, v := range wasFound {
		res = append(res, goodComics{id: k, numOfKeywords: v})
	}
	slices.SortFunc(res, func(a, b goodComics) int {
		return cmp.Compare(a.numOfKeywords, b.numOfKeywords) * (-1)
	})
	return res
}

type goodComics struct {
	id            int
	numOfKeywords int
}
