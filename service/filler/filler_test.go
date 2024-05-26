package filler

import (
	"context"
	"courses/core"
	"courses/service/xkcd"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
)

type db struct{}

func (d db) Read() (map[int]core.ComicsDescript, error) {
	return make(map[int]core.ComicsDescript), nil
}

func (d db) Write(core.ComicsDescript, int) error {
	return nil
}

func TestFiller_FillAllComics(t *testing.T) {
	dwnl := xkcd.NewComicsDownloader("https://xkcd.com")
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	testSet := map[int]core.ComicsDescript{1: {Keywords: []string{"banana", "minion"}, Url: "xkcd.com"}}

	filler := NewFiller(100, testSet, db{}, dwnl, *logger)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()
	comics, err := filler.FillMissedComics(ctx)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fmt.Println(len(comics))
	if len(comics) == 0 {
		t.Errorf("comics didn't downloaded")
	}
}
