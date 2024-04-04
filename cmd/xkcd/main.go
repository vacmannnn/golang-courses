package main

import (
	"courses/pkg/database"
	"courses/pkg/xkcd"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	SourceUrl string `yaml:"source_url"`
	DBFile    string `yaml:"db_file"`
}

func newConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	var numOfComics int
	flag.IntVar(&numOfComics, "n", -1, "number of comics to save")
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config.yml file")
	var showDownloadedComics bool
	flag.BoolVar(&showDownloadedComics, "o", false, "show info about downloaded comics")
	flag.Parse()

	conf, err := newConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	bytes, err := xkcd.GetComicsFromSite(conf.SourceUrl, conf.DBFile)
	if err != nil {
		log.Println(err)
		if bytes == nil {
			return
		}
	}

	if bytes != nil && showDownloadedComics {
		var comicsToJSON map[int]xkcd.ComicsDescript
		err = json.Unmarshal(bytes, &comicsToJSON)
		if err != nil {
			log.Println(err)
			return
		}
		if numOfComics != -1 {
			for i, comics := range comicsToJSON {
				fmt.Printf("id - %d, keywords - %s, url - %s\n", i, comics.Keywords, comics.Url)
			}
		} else {
			for i := 1; i <= numOfComics; i++ {
				fmt.Printf("id - %d, keywords - %s, url - %s\n", i, comicsToJSON[i].Keywords,
					comicsToJSON[i].Url)
			}
		}
	}

	if err = database.WriteToDB(conf.DBFile, bytes); err != nil {
		log.Fatal(err)
	}
}
