package catalog

import (
	"cmp"
	"courses/internal/core"
	"slices"
)

type goodComics struct {
	Id            int
	NumOfKeywords int
}

// FindByIndex searches input string by its slice of keywords and returns slice of most suitable comics URLs. The more
// comics suitable, the lower the index
func (c *ComicsCatalog) FindByIndex(input []string) []string {
	wasFound := make(map[int]int)
	for _, keywords := range input {
		for _, comicsID := range c.index[keywords] {
			wasFound[comicsID]++
		}
	}
	var res []goodComics
	for k, v := range wasFound {
		if v != 0 {
			res = append(res, goodComics{Id: k, NumOfKeywords: v})
		}
	}
	slices.SortFunc(res, func(a, b goodComics) int {
		return cmp.Compare(a.NumOfKeywords, b.NumOfKeywords) * (-1)
	})
	var urls []string
	for i := 0; i < min(core.MaxComicsToShow, len(res)); i++ {
		urls = append(urls, c.comics[res[i].Id].Url)
	}
	return urls
}

// findByComics unused cause of inefficient speed compared to FindByIndex
func (c *ComicsCatalog) findByComics(input []string) []string {
	var res []goodComics
	for id, v := range c.comics {
		var numOfWords int
		for _, word := range input {
			if slices.Contains(v.Keywords, word) {
				numOfWords++
			}
		}
		if numOfWords != 0 {
			res = append(res, goodComics{Id: id, NumOfKeywords: numOfWords})
		}
	}
	slices.SortFunc(res, func(a, b goodComics) int {
		return cmp.Compare(a.NumOfKeywords, b.NumOfKeywords) * (-1)
	})
	var urls []string
	for i := 0; i < min(core.MaxComicsToShow, len(res)); i++ {
		urls = append(urls, c.comics[res[i].Id].Url)
	}
	return urls
}
