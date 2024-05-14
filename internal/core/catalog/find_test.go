package catalog

import (
	"context"
	"courses/internal/core"
	"courses/internal/core/filler"
	"courses/internal/core/xkcd"
	"courses/internal/database"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkDiffMethToSearch(b *testing.B) {
	myDB, _ := database.NewDB("test.json")

	comics, _ := myDB.Read()
	if comics == nil {
		comics = make(map[int]core.ComicsDescript, 3000)
	}

	sourceUrl := "xkcd.com"
	downloader := xkcd.NewComicsDownloader(sourceUrl)

	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := slog.New(handler)
	comicsFiller := filler.NewFiller(core.GoroutineNum, comics, myDB, downloader, *logger)
	comics, _ = comicsFiller.FillMissedComics(context.Background())

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
	testString := []string{"my favorite comics is about unknown mystery person", "idk what comics to search",
		"cool banana man", "orange box sits under that orange table and takes orange to make orange juice",
		"funny comics about math"}
	for _, str := range testString {
		comicsName := "findFindByIndex-" + strconv.Itoa(len(str))
		catalog := NewCatalog(comics, comicsFiller)
		b.Run(comicsName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				catalog.FindByIndex(strings.Split(str, " "))
			}
		})
		comicsName = "findByComics-" + strconv.Itoa(len(str))
		b.Run(comicsName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				catalog.findByComics(strings.Split(str, " "))
			}
		})
	}
}
