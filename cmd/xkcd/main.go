package main

import (
	"context"
	"courses/internal/adapter/handler"
	"courses/internal/core"
	"courses/internal/core/catalog"
	"courses/internal/core/filler"
	"courses/internal/core/xkcd"
	"courses/internal/database"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	myDB, err := database.NewDB(conf.DBFile)
	if err != nil {
		log.Fatal(err)
	}
	defer myDB.Close()

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
	ctlg := catalog.NewCatalog(comics, comicsFiller)

	mux := handler.NewMux(ctlg, logger)
	portStr := fmt.Sprintf(":%d", port)

	// based on https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
	server := &http.Server{Addr: portStr, Handler: mux}
	go func() {
		if err = server.ListenAndServe(); err != nil {
			logger.Debug("server error", "err", err.Error())
		}
	}()
	logger.Info("Server started")

	c := cron.New()
	_, err = c.AddFunc("0 13 * * *", func() {
		logger.Info("Send update")
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		url := fmt.Sprintf("http://localhost:%d/update", port)
		req, err := http.NewRequestWithContext(context.Background(),
			http.MethodPost, url, nil)
		if err != nil {
			logger.Error(err.Error())
		}

		res, err := client.Do(req)
		if err != nil {
			logger.Error(err.Error())
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(res.Body)
		if res.StatusCode != http.StatusOK {
			logger.Error(fmt.Sprintf("unexpected status: got %v", res.Status))
		}
	})
	if err != nil {
		logger.Error("Cron error", "err", err.Error())
	}
	c.Start()

	<-ctx.Done()

	ctx, stop = context.WithTimeout(context.Background(), core.MaxWaitTime)
	defer stop()
	if err = server.Shutdown(ctx); err != nil {
		logger.Debug("server shutdown error", "err", err.Error())
	}
}
