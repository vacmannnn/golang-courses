package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"slices"
	"strings"
	"testing"
)

func BenchmarkFindByIndex(b *testing.B) {
	conf, _ := newConfig("../../config.yaml")

	myDB := database.NewDB(conf.DBFile)

	comics, _ := myDB.Read()
	if comics == nil {
		comics = make(map[int]core.ComicsDescript, 3000)
	}
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)

	comics, _ = fillMissedComics(5, comics, myDB, downloader)

	index := make(map[string][]int)
	var doc []string
	for k, v := range comics {
		doc = slices.Concat(doc, v.Keywords)
		for i, token := range v.Keywords {
			if !slices.Contains(v.Keywords[:i], token) {
				index[token] = append(index[token], k)
			}
		}
	}
	b.Run("findByIndex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			findByIndex(index, strings.Split("my favorite comics is about unknown mystery person", " "))
		}
	})
	b.Run("findByComics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			findByComics(comics, strings.Split("my favorite comics is about unknown mystery person", " "))
		}
	})
}
