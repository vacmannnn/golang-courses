package xkcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

// GetComicsFromSite gets url, id of first comics to download and last. If any value is not greater than 0
// it will be reassigned to 1 in case of first comics and to latest comics at whole cite in case of last id.
// Function will log any non-critical error.
func GetComicsFromSite(urlName string, startComicsId, endComicsId int) (map[int]ComicsDescript, error) {
	if endComicsId < 1 {
		latestComics, err := getComicsFromURL("https://xkcd.com/info.0.json")
		if err != nil {
			return nil, err
		}
		endComicsId = latestComics.Num
	}
	if startComicsId < 1 {
		startComicsId = 1
	}
	if endComicsId < startComicsId {
		return nil, errors.New("id of last comics to download should be lower than first comics to download")
	}

	var curGoroutines int
	var wg sync.WaitGroup
	var comicsMutex sync.Mutex
	comicsToJSON := make(map[int]ComicsDescript, endComicsId-startComicsId)

	for i := startComicsId; i <= endComicsId; i++ {
		wg.Add(1)
		curGoroutines++
		go func(comicsID int) {
			comicsURL := fmt.Sprintf("%s/%d/info.0.json", urlName, comicsID)
			log.Println(comicsURL)

			myComics, err := getComicsFromURL(comicsURL)
			if err != nil {
				log.Printf("%s, comicsID is - %d", err, comicsID)
			}

			keywords := strings.Split(myComics.Transcript, " ")
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

	return comicsToJSON, nil
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
