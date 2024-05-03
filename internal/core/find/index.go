package catalog

import (
	"context"
	"courses/internal/core"
	"courses/internal/core/filler"
	"reflect"
	"slices"
	"sync"
)

// как-то грустно, что finder внутри себя хранит доступ к БД (filler)

type ComicsCatalog struct {
	comics map[int]core.ComicsDescript
	index  map[string][]int
	filler filler.Filler
	mt     sync.Mutex
}

func NewCatalog(comics map[int]core.ComicsDescript, filler filler.Filler) *ComicsCatalog {
	f := ComicsCatalog{comics: comics, filler: filler}
	f.buildIndex()
	return &f
}

func (f *ComicsCatalog) buildIndex() {
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

func (f *ComicsCatalog) GetIndex() map[string][]int {
	return f.index
}

func (f *ComicsCatalog) UpdateComics() (map[string]int, error) {
	updatedComics, err := f.filler.FillMissedComics(context.Background())
	if err != nil {
		return nil, err
	}

	eq := reflect.DeepEqual(updatedComics, f.comics)
	var n int
	if !eq {
		for k, v := range updatedComics {
			if slices.Equal(f.comics[k].Keywords, v.Keywords) {
				n++
			}
		}

		// need to update current comics set with corresponding index
		f.mt.Lock()
		f.comics = updatedComics
		f.buildIndex()
		f.mt.Unlock()
	}

	diff := map[string]int{
		"new": n, "total": len(updatedComics),
	}
	return diff, nil
}
