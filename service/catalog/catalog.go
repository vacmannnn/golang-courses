package catalog

import (
	"context"
	"courses/core"
	"courses/service/filler"
	"maps"
	"reflect"
	"slices"
	"sync"
)

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

func (c *ComicsCatalog) buildIndex() {
	index := make(map[string][]int)
	for k, v := range c.comics {
		for i, token := range v.Keywords {
			if !slices.Contains(v.Keywords[:i], token) {
				index[token] = append(index[token], k)
			}
		}
	}
	c.index = index
}

func (c *ComicsCatalog) GetIndex() map[string][]int {
	return c.index
}

func (c *ComicsCatalog) UpdateComics() (int, int, error) {
	oldComics := make(map[int]core.ComicsDescript)
	maps.Copy(oldComics, (*c).comics)

	updatedComics, err := c.filler.FillMissedComics(context.Background())
	if err != nil {
		return 0, 0, err
	}

	eq := reflect.DeepEqual(updatedComics, oldComics)
	var n int
	if !eq {
		for k, v := range updatedComics {
			if !slices.Equal(oldComics[k].Keywords, v.Keywords) {
				n++
			}
		}

		// need to update current comics set with corresponding index
		c.mt.Lock()
		c.comics = updatedComics
		c.buildIndex()
		c.mt.Unlock()
	}

	return n, len(updatedComics), nil
}
