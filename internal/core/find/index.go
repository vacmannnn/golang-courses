package find

import (
	"courses/internal/core"
	"slices"
)

type Finder struct {
	comics map[int]core.ComicsDescript
	index  map[string][]int
}

type goodComics struct {
	Id            int
	NumOfKeywords int
}

func NewFinder(comics map[int]core.ComicsDescript) Finder {
	f := Finder{comics: comics}
	f.buildIndex()
	return f
}

func (f *Finder) buildIndex() {
	index := make(map[string][]int)
	for k, v := range f.comics {
		for i, token := range v.Keywords {
			if !slices.Contains(v.Keywords[:i], token) {
				index[token] = append(index[token], k)
			}
		}
	}
	f.index = index
}

func (f *Finder) GetIndex() map[string][]int {
	return f.index
}
