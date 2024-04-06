package main

import (
	"courses/pkg/database"
	"courses/pkg/words"
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
	var numOfComics int
	flag.IntVar(&numOfComics, "n", -1, "number of comics to save")
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "path to config.yml file")
	var showDownloadedComics bool
	flag.BoolVar(&showDownloadedComics, "o", false, "show info about downloaded comics")
	flag.Parse()

	// get config
	conf, err := newConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// read existed json to simplify downloading
	comicsToJSON := make(map[int]xkcd.ComicsDescript)
	myDB := database.NewDB(conf.DBFile)

	// it's ok if there was an error in file because we are going to create again and overwrite it
	file, err := myDB.Read()
	if err != nil {
		log.Println(err)
	}

	// if case of any error in unmarshalling whole file will be overwritten due to corruption
	err = json.Unmarshal(file, &comicsToJSON)
	lastComicsNum := 1
	if err != nil {
		log.Println(err)
		lastComicsNum = 1
	} else {
		for k := range comicsToJSON {
			lastComicsNum = max(lastComicsNum, k)
		}
	}

	// download needed
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl)
	comics, err := downloader.GetComicsFromSite(lastComicsNum, numOfComics)
	if err != nil {
		log.Println(err)
		if comics == nil {
			return
		}
	}
	for k, v := range comics {
		v.Keywords = words.StemStringWithClearing(v.Keywords)
		comicsToJSON[k] = v
	}

	// show if needed
	if showDownloadedComics {
		if numOfComics == -1 {
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

	// load to JSON
	bytes, err := marshallComics(comicsToJSON)
	if err != nil {
		log.Println(err)
		return
	}

	if err = myDB.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

func marshallComics(comicsToJSON map[int]xkcd.ComicsDescript) ([]byte, error) {
	return json.MarshalIndent(comicsToJSON, "", " ")
}
