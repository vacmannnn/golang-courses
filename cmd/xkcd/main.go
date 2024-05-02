package main

import (
	"courses/internal/adapter/handler"
	"courses/internal/core"
	"courses/internal/core/find"
	"courses/internal/core/xkcd"
	"courses/internal/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	configPath, port, loggerLevel := getFlags()

	opts := &slog.HandlerOptions{
		Level: loggerLevel,
	}
	logHandler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(logHandler)

	conf, err := getConfig(configPath)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	// -1 is default port value if there is no -p flag
	if port == -1 {
		port = conf.Port
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

	filler := xkcd.NewFiller(goroutineNum, comics, myDB, downloader, *logger)
	comics, err = filler.FillMissedComics()
	if err != nil {
		logger.Error(err.Error())
	}

	// build index
	finder := find.NewFinder(comics, filler)
	index := finder.GetIndex()

	// write to index.json
	file, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		logger.Warn(err.Error())
	}

	err = os.WriteFile("index.json", file, 0644)
	if err != nil {
		logger.Warn(err.Error())
	}

	mux := handler.CreateServeMux(finder)
	portStr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portStr, mux)
}
