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
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

type comicsDescriptWithID struct {
	core.ComicsDescript
	id int
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

	// read existed json to simplify downloading
	myDB := database.NewDB(conf.DBFile)

	// it's ok if there was an error in file because we are going to create again and overwrite it
	comicsToJSON, err := myDB.Read()
	if comicsToJSON == nil {
		comicsToJSON = make(map[int]core.ComicsDescript, 3000)
	}
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d comics in base", len(comicsToJSON))

	// comics downloads by parts. Each parts consist of N (current num is 500) comics, after downloading each part
	// it will be uploaded to DB to prevent problems with unexpected program kill
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)

	comicsIDChan := make(chan int)
	comicsChan := make(chan comicsDescriptWithID)

	go func() {
		for curID := range comicsIDChan {
			if comicsToJSON[curID].Keywords == nil {
				descript, id, err := downloader.GetComicsFromID(curID)
				if err != nil {
					comicsChan <- comicsDescriptWithID{}
					close(comicsChan)
					return
				}
				comicsChan <- comicsDescriptWithID{id: id, ComicsDescript: descript}
				continue
			}
			comicsChan <- comicsDescriptWithID{id: curID, ComicsDescript: comicsToJSON[curID]}
		}
	}()

	for i := 1; ; i++ {
		var curComics comicsDescriptWithID
		comicsIDChan <- i
		curComics = <-comicsChan
		if curComics.Url != "" {
			curComics.Keywords = words.StemStringWithClearing(curComics.Keywords)
			var comics = make(map[int]core.ComicsDescript)
			comics[curComics.id] = curComics.ComicsDescript
			if err = myDB.Write(comics); err != nil {
				log.Fatal(err)
			}
		} else {
			close(comicsIDChan)
			for curComics = range comicsChan {
				if curComics.Url != "" {
					curComics.Keywords = words.StemStringWithClearing(curComics.Keywords)
					comics := map[int]core.ComicsDescript{curComics.id: curComics.ComicsDescript}
					if err = myDB.Write(comics); err != nil {
						log.Fatal(err)
					}
				}
			}
			break
		}
	}
}
