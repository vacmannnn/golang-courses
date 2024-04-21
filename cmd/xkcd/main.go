package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
)

type comicsDescriptWithID struct {
	core.ComicsDescript
	id int
}

func main() {
	// parse flags
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "path to config.yml file")
	var inputString string
	flag.StringVar(&inputString, "s", "", "string to find")
	var byIndex bool
	flag.BoolVar(&byIndex, "i", false, "find comics by index")
	flag.Parse()
	if inputString == "" {
		// TODO
		log.Println("empty input")
	}
	clearedInput := words.StemStringWithClearing(strings.Split(inputString, " "))

	// get config
	conf, err := newConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	goroutineNum, err := getGoroutinesNum()
	if err != nil {
		log.Println(err)
	}

	// read existed DB to simplify downloading
	myDB := database.NewDB(conf.DBFile)

	comics, err := myDB.Read()
	if comics == nil {
		comics = make(map[int]core.ComicsDescript, 3000)
	}
	if err != nil {
		log.Println(err)
	}
	log.Printf("%d comics in base", len(comics))
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)

	comics, err = fillMissedComics(goroutineNum, comics, myDB, downloader)
	if err != nil {
		log.Println(err)
	}

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
	var res []goodComics
	if byIndex {
		res = findByIndex(index, clearedInput)
	} else {
		res = findByComics(comics, clearedInput)
	}
	for i := 0; i < min(10, len(res)); i++ {
		fmt.Println(res[i], comics[res[i].id].Url)
	}
}
