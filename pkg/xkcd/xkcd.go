package xkcd

import (
	"encoding/json"
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

type ComicsDownloader struct {
	comicsURL string
}

// NewComicsDownloader sets link to source cite with comics
func NewComicsDownloader(comicsURL string) ComicsDownloader {
	return ComicsDownloader{comicsURL: comicsURL}
}

// GetComicsFromSite gets id of first comics to download and last. If any value is not greater than 0
// it will be reassigned to 1 in case of first comics and to latest comics at whole cite in case of last id.
// Function will log any non-critical error.
func (c ComicsDownloader) GetComicsFromSite(comicsID []int) (map[int]ComicsDescript, error) {
	var latestComicsID int
	if len(comicsID) == 0 {
		comicsURL := fmt.Sprintf("%s/info.0.json", c.comicsURL)
		latestComics, err := c.getComicsFromURL(comicsURL)
		if err != nil {
			return nil, err
		}
		latestComicsID = latestComics.Num
		for i := 1; i <= latestComicsID; i++ {
			comicsID = append(comicsID, i)
		}
	}

	var curGoroutines int
	var wg sync.WaitGroup
	comicsChan := make(chan comicsInfo)
	comicsToJSON := make(map[int]ComicsDescript, len(comicsID))

	go func() {
		for _, i := range comicsID {
			wg.Add(1)
			curGoroutines++
			go func(comicsID int) {
				comicsURL := fmt.Sprintf("%s/%d/info.0.json", c.comicsURL, comicsID)
				log.Println(comicsURL)

				myComics, err := c.getComicsFromURL(comicsURL)
				if err != nil {
					log.Printf("%s, comicsID is - %d", err, comicsID)
				}

				comicsChan <- myComics
				wg.Done()
			}(i)

			// Need to download step by step due possible heavy load on the network
			if curGoroutines%500 == 0 {
				wg.Wait()
			}
		}
	}()
	for range comicsID {
		comicsOwner := <-comicsChan
		keywords := strings.Split(comicsOwner.Transcript, " ")
		comicsToJSON[comicsOwner.Num] = ComicsDescript{Url: comicsOwner.ImgURL, Keywords: keywords}
	}
	wg.Wait()

	return comicsToJSON, nil
}

func (c ComicsDownloader) getComicsFromURL(comicsURL string) (comicsInfo, error) {
	client := http.Client{}
	resp, err := client.Get(comicsURL)
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
