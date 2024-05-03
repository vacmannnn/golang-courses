package core

import "time"

type ComicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type ComicsDownloader interface {
	GetComicsFromID(int) (ComicsDescript, int, error)
}

const MaxWaitTime = time.Second * 5

const MaxComicsToShow = 10

const GoroutineNum = 250
