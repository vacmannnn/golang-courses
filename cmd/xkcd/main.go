package main

import (
	"context"
	"courses/internal/adapter/handler"
	"courses/internal/core"
	"courses/internal/core/filler"
	"courses/internal/core/find"
	"courses/internal/core/xkcd"
	"courses/internal/database"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)
	comicsFiller := filler.NewFiller(core.GoroutineNum, comics, myDB, downloader, *logger)
	comics, err = comicsFiller.FillMissedComics(ctx)
	if err != nil {
		logger.Error(err.Error())
	}

	// build index
	finder := find.NewFinder(comics, comicsFiller)
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

	mux := handler.CreateServeMux(finder, logger)
	portStr := fmt.Sprintf(":%d", port)

	// based on https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
	server := &http.Server{Addr: portStr, Handler: mux}
	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Debug("server error", "err", err.Error())
		}
	}()
	<-ctx.Done()

	ctx, stop = context.WithTimeout(context.Background(), core.MaxWaitTime)
	defer stop()
	if err = server.Shutdown(ctx); err != nil {
		logger.Debug("server shutdown error", "err", err.Error())
	}
}
