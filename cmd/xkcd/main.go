package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"flag"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

type comicsDescriptWithID struct {
	core.ComicsDescript
	id int
}

func main() {
	// parse flags
	var configPath string
	flag.StringVar(&configPath, "c", "config.yaml", "path to config.yml file")
	flag.Parse()

	// get config
	conf, err := newConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}
	goroutineNum := 500
	// if err != nil {
	// 	log.Println(err)
	// }

	// read existed DB to simplify downloading
	myDB := database.NewDB(conf.DBFile)

	comicsToJSON, err := myDB.Read()
	if comicsToJSON == nil {
		comicsToJSON = make(map[int]core.ComicsDescript, 3000)
	}
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d comics in base", len(comicsToJSON))

	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)

	comicsIDChan := make(chan int, goroutineNum)
	comicsChan := make(chan comicsDescriptWithID, goroutineNum)

	for range goroutineNum {
		go worker(downloader, comicsToJSON, comicsIDChan, comicsChan)
	}

	var curComics comicsDescriptWithID
	for i := 1; ; i++ {
		// send in advance bunch of ID to optimize downloading
		if i%goroutineNum == 1 {
			for j := i; j < i+goroutineNum; j++ {
				comicsIDChan <- j
			}
		}

		curComics = <-comicsChan
		if curComics.Url != "" {
			if err = writeComicsWithID(curComics, &myDB); err != nil {
				log.Fatal(err)
			}
		} else {
			// wait till we get all comics from site
			time.Sleep(core.MaxWaitTime)
			for range len(comicsChan) {
				curComics = <-comicsChan
				if curComics.Url == "" {
					continue
				}
				if err = writeComicsWithID(curComics, &myDB); err != nil {
					log.Fatal(err)
				}
			}
			close(comicsIDChan)
			close(comicsChan)
			break
		}
	}
}

func worker(downloader xkcd.ComicsDownloader, comics map[int]core.ComicsDescript, comicsIDChan <-chan int,
	results chan<- comicsDescriptWithID) {
	for comID := range comicsIDChan {
		if comics[comID].Keywords == nil {
			descript, id, err := downloader.GetComicsFromID(comID)
			if err != nil {
				results <- comicsDescriptWithID{id: comID}
				continue
			}
			descript.Keywords = words.StemStringWithClearing(descript.Keywords)
			results <- comicsDescriptWithID{id: id, ComicsDescript: descript}
			continue
		}
		results <- comicsDescriptWithID{id: comID, ComicsDescript: comics[comID]}
	}
}

func writeComicsWithID(comicsWID comicsDescriptWithID, db *database.DataBase) error {
	var comics = make(map[int]core.ComicsDescript)
	comics[comicsWID.id] = comicsWID.ComicsDescript
	return db.Write(comics)
}

func newConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func getGoroutinesNum() (int, error) {
	defaultValue := 500
	file, err := os.Open("parallel")
	if err != nil {
		return defaultValue, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)
	if err = d.Decode(&defaultValue); err != nil {
		return 500, err
	}

	return defaultValue, nil
}
