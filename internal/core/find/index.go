package find

import (
	"courses/internal/core"
	"courses/internal/core/xkcd"
	"reflect"
	"slices"
)

// TODO: подумать над set new filler
// как-то грустно, что filler внутри себя хранит доступ к БД

type Finder struct {
	comics map[int]core.ComicsDescript
	filler xkcd.Filler
	index  map[string][]int
}

type goodComics struct {
	Id            int
	NumOfKeywords int
}

func NewFinder(comics map[int]core.ComicsDescript, filler xkcd.Filler) Finder {
	f := Finder{comics: comics, filler: filler}
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

func (f *Finder) UpdateComics() map[string]int {
	updatedComics, err := f.filler.FillMissedComics()
	if err != nil {
		// TODO
	}
	eq := reflect.DeepEqual(updatedComics, f.comics)
	var n int
	if !eq {
		for k, v := range updatedComics {
			if slices.Equal(f.comics[k].Keywords, v.Keywords) {
				n++
			}
		}
		// TODO: shared memory, case with everyday update
		f.comics = updatedComics
	}
	diff := map[string]int{
		"new": n, "total": len(updatedComics),
	}
	return diff
}
