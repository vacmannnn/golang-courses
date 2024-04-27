package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type comicsDescriptWithID struct {
	core.ComicsDescript
	id int
}

func main() {
	configPath, inputString, byIndex, loggerLevel := getFlags()

	opts := &slog.HandlerOptions{
		Level: loggerLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	if inputString == "" {
		logger.Warn("Input string shouldn't be empty")
	}
	clearedInput := words.StemStringWithClearing(strings.Split(inputString, " "))

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

	// find comics
	var res []goodComics
	if byIndex {
		res = findByIndex(index, clearedInput)
	} else {
		res = findByComics(comics, clearedInput)
	}

	const maxComicsToShow = 10
	for i := 0; i < min(maxComicsToShow, len(res)); i++ {
		fmt.Println(res[i], comics[res[i].id].Url)
	}
}
