package main

import (
	"courses/pkg/database"
	"courses/pkg/words"
	"courses/pkg/xkcd"
	"flag"
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
		comicsToJSON = make(map[int]xkcd.ComicsDescript, 3000)
	}
	if err != nil {
		log.Println(err)
	}

	log.Printf("%d comics in base", len(comicsToJSON))

	// comics downloads by parts. Each parts consist of N (current num is 500) comics, after downloading each part
	// it will be uploaded to DB to prevent problems with unexpected program kill
	downloader := xkcd.NewComicsDownloader(conf.SourceUrl, comicsToJSON)

	var comics map[int]xkcd.ComicsDescript
	for comicsToDownload := 500; comicsToDownload == 500; {
		comics, comicsToDownload, err = downloader.GetNComicsFromSite(comicsToDownload)
		if err != nil {
			log.Println(err)
		}
		for k, v := range comics {
			v.Keywords = words.StemStringWithClearing(v.Keywords)
			comicsToJSON[k] = v
		}
		if err = myDB.Write(comicsToJSON); err != nil {
			log.Fatal(err)
		}
	}
}
