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
	"time"
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
	comicsToJSON := make(map[int]xkcd.ComicsDescript)
	myDB := database.NewDB(conf.DBFile)

	// it's ok if there was an error in file because we are going to create again and overwrite it
	file, err := myDB.Read()
	if err != nil {
		log.Println(err)
	}

	// if case of any error in unmarshalling whole file will be overwritten due to corruption
	err = json.Unmarshal(file, &comicsToJSON)
	fmt.Println(comicsToJSON)

	// download needed
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl, comicsToJSON)

	var comics map[int]xkcd.ComicsDescript
	for comicsToDownload := 500; comicsToDownload == 500; {
		time.Sleep(time.Second * 3)
		comics, comicsToDownload, err = downloader.GetComicsFromSite(comicsToDownload)
		fmt.Println("all good!", comicsToDownload)
	}
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
