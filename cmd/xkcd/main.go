package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"
)

type comicsDescriptWithID struct {
	core.ComicsDescript
	id int
}

func main() {
	configPath, _, _, loggerLevel := getFlags()

	opts := &slog.HandlerOptions{
		Level: loggerLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	conf, err := getConfig(configPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	goroutineNum, err := getGoroutinesNum()
	if err != nil {
		logger.Error(err.Error())
	}

	// read existed DB to simplify downloading
	myDB := database.NewDB(conf.DBFile)

	comics, err := myDB.Read()
	if comics == nil {
		comics = make(map[int]core.ComicsDescript, 3000)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("base opened", "comics in base", len(comics))
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)

	f := newFiller(goroutineNum, comics, myDB, downloader, *logger)
	comics, err = f.fillMissedComics()
	if err != nil {
		logger.Error(err.Error())
	}

	// build index
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

	// write to index.json
	file, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		logger.Warn(err.Error())
	}

	err = os.WriteFile("index.json", file, 0644)
	if err != nil {
		logger.Warn(err.Error())
	}

	const maxComicsToShow = 10
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pics", func(wr http.ResponseWriter, r *http.Request) {
		comicsKeywords := r.URL.Query().Get("search")
		clearedKeywords := words.StemStringWithClearing(strings.Split(comicsKeywords, " "))
		res := findByIndex(index, clearedKeywords)
		var urls []string
		for i := 0; i < min(maxComicsToShow, len(res)); i++ {
			urls = append(urls, comics[res[i].id].Url)
		}
		data, _ := json.Marshal(urls)
		_, _ = wr.Write(data)
	})
	mux.HandleFunc("POST /update", func(wr http.ResponseWriter, r *http.Request) {
		updatedComics, err := f.fillMissedComics()
		if err != nil {
			// TODO
		}
		eq := reflect.DeepEqual(updatedComics, comics)
		var data []byte
		var n int
		if !eq {
			for k, v := range updatedComics {
				if slices.Equal(comics[k].Keywords, v.Keywords) {
					n++
				}
			}
			// TODO: shared memory, case with everyday update
			comics = updatedComics
		}
		diff := map[string]int{
			"new": n, "total": len(updatedComics),
		}
		data, err = json.Marshal(diff)
		if err != nil {
			// TODO
		}
		wr.Write(data)
	})
	http.ListenAndServe(":8080", mux)
}
