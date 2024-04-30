package find

import (
	"cmp"
	"courses/internal/core"
	"slices"
)

// TODO: struct `finder`
// Было бы круто инкапсулировать индекс, просто сделать структурку, которая получает комиксы и сама строит индекс
// Комикс можно было бы заапдейтить

// И переименовать было бы неплохо

func (f *Finder) ByIndex(input []string) []string {
	wasFound := make(map[int]int)
	for _, keywords := range input {
		for _, comicsID := range f.index[keywords] {
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
		urls = append(urls, f.comics[res[i].Id].Url)
	}
	return urls
}

// ByComics unused cause of inefficient speed compared to ByIndex
func (f *Finder) byComics(input []string) []string {
	var res []goodComics
	for id, v := range f.comics {
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
		urls = append(urls, f.comics[res[i].Id].Url)
	}
	return urls
}
