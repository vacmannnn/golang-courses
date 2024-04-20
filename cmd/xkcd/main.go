package main

import (
	"cmp"
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"slices"
	"strings"
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
	var inputString string
	flag.StringVar(&inputString, "s", "", "string to find")
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
	for k, v := range comics {
		for _, token := range v.Keywords {
			index[token] = append(index[token], k)
		}
	}
	res := findByIndex(index, clearedInput)
	for i := 0; i < min(len(res), 10); i++ {
		fmt.Println(res[i], comics[res[i].id].Url)
	}
}

func findByIndex(index map[string][]int, input []string) []goodComics {
	wasFound := make(map[int]int)
	for _, keywords := range input {
		for _, comicsID := range index[keywords] {
			wasFound[comicsID]++
		}
	}
	var res []goodComics
	for k, v := range wasFound {
		res = append(res, goodComics{id: k, numOfKeywords: v})
	}
	slices.SortFunc(res, func(a, b goodComics) int {
		return cmp.Compare(a.numOfKeywords, b.numOfKeywords) * (-1)
	})
	return res
}

type goodComics struct {
	id            int
	numOfKeywords int
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
	obj := make(map[string]int)

	yamlFile, err := os.ReadFile("parallel")
	if err != nil {
		return defaultValue, err
	}
	err = yaml.Unmarshal(yamlFile, obj)
	if err != nil {
		return defaultValue, err
	}

	if obj["goroutines"] == 0 {
		obj["goroutines"] = defaultValue
	}
	return obj["goroutines"], nil
}
