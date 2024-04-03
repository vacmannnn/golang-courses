package main

import (
	"courses/pkg/database"
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
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err = d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// cross-platforming ? will it work well in each platform ?
	conf, err := newConfig("../../config.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}
	// fmt.Println(conf.SourceUrl, conf.DBFile, err)
	var numOfComics int
	flag.IntVar(&numOfComics, "n", 1, "number of comics to save")
	flag.Parse()
	// TODO: return error
	bytes := xkcd.GetNComicsFromSite(conf.SourceUrl, conf.DBFile, numOfComics)
	database.WriteToDB(conf.DBFile, bytes)
}
