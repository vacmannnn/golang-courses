package core

import "time"

type Config struct {
	SourceUrl        string `yaml:"source_url"`
	DBFile           string `yaml:"db_file"`
	Port             int    `yaml:"port"`
	ConcurrencyLimit int    `yaml:"concurrency_limit"`
	RateLimit        int    `yaml:"rate_limit"`
	TokenMaxTime     int    `yaml:"token_max_time"`
}

type ComicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type ComicsDownloader interface {
	GetComicsFromID(int) (ComicsDescript, int, error)
}

type DataBase interface {
	Write(ComicsDescript, int) error
	Read() (map[int]ComicsDescript, error)
}

type Catalog interface {
	FindByIndex([]string) []string
	GetIndex() map[string][]int
	UpdateComics() (int, int, error)
}

const MaxWaitTime = time.Second * 5

const MaxComicsToShow = 10

const GoroutineNum = 250
