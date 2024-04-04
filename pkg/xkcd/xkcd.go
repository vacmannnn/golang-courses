package xkcd

import (
	"courses/pkg/database"
	"courses/pkg/words"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type comicsInfo struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	ImgURL     string `json:"img"`
	Title      string `json:"string"`
	Day        string `json:"day"`
}

type ComicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// GetNComicsFromSite gets url, name of existing DB file and number of comics to download. If db file doesn't exist it
// possible to pass "". Num of comics should be greater than 0. Function will log any non-critical error.
// Returned slice of byte may be not nil if some comics downloaded.
func GetNComicsFromSite(urlName string, dbFileName string, comicsNum int) ([]byte, error) {
	if comicsNum < 1 {
		return nil, errors.New("number of comics should be greater than 0, default value is 1")
	}
	comicsToJSON := make(map[int]ComicsDescript, comicsNum)
	comicsMutex := sync.RWMutex{}

	// it's ok if there was an error in file because we are going to create again and overwrite it
	file, err := database.ReadFromDB(dbFileName)
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

	wg := sync.WaitGroup{}
	var curGoroutines int
	// last comicsInfo will be overwritten due possible corruption
	for i := lastComicsNum; i <= comicsNum; i++ {
		wg.Add(1)
		curGoroutines++
		go func(comicsID int) {
			comicsURL := fmt.Sprintf("%s/%d/info.0.json", urlName, comicsID)
			log.Println(comicsURL)

			myComics, err := getComicsFromURL(comicsURL)
			if err != nil {
				log.Printf("%s, comicsID is - %d", err, comicsID)
			}

			keywords := words.StemStringWithClearing(myComics.Transcript)
			comicsMutex.Lock()
			comicsToJSON[comicsID] = ComicsDescript{Url: myComics.ImgURL, Keywords: keywords}
			comicsMutex.Unlock()
			wg.Done()
		}(i)

		// Need to download step by step due possible heavy load on the network
		if curGoroutines%500 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()
	bytes, err := marshallComics(comicsToJSON)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func getComicsFromURL(comicsURL string) (comicsInfo, error) {
	c := http.Client{}
	resp, err := c.Get(comicsURL)
	if err != nil {
		return comicsInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return comicsInfo{}, err
	}

	var myComics comicsInfo
	err = json.Unmarshal(body, &myComics)
	if err != nil {
		return comicsInfo{}, err
	}

	return myComics, nil
}

func marshallComics(comicsToJSON map[int]ComicsDescript) ([]byte, error) {
	bytes, err := json.MarshalIndent(comicsToJSON, "", " ")
	return bytes, err
}
