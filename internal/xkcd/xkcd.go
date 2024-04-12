package xkcd

import (
	"courses/internal/core"
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

type ComicsDownloader struct {
	comicsURL        string
	comics           map[int]core.ComicsDescript
	lastDownloadedID int
}

// NewComicsDownloader sets link to source site with comics
func NewComicsDownloader(comicsURL string, comics map[int]core.ComicsDescript) ComicsDownloader {
	return ComicsDownloader{comicsURL: comicsURL, comics: comics}
}

// ChangeStartIDOfDownloading sets new ID from which newest comics will be downloaded, default value will
// grow after each download by numOfComics
func (c *ComicsDownloader) ChangeStartIDOfDownloading(startPos int) {
	c.lastDownloadedID = startPos
}

// GetNComicsFromSite gets num of comics that will be downloaded from site. Start id of downloading comics depends
// on previous downloads. Returns map with needed comics and number of successful downloads.
func (c *ComicsDownloader) GetNComicsFromSite(numOfComics int) (map[int]core.ComicsDescript, int, error) {

	var resMap = make(map[int]core.ComicsDescript, numOfComics)
	var (
		wg                sync.WaitGroup
		mt                sync.Mutex
		successDownloaded int
	)

	for i := c.lastDownloadedID + 1; i <= c.lastDownloadedID+numOfComics; i++ {
		wg.Add(1)
		go func(comicsID int) {
			defer wg.Done()
			if c.comics[comicsID].Keywords != nil {
				mt.Lock()
				defer mt.Unlock()
				resMap[comicsID] = c.comics[comicsID]
				successDownloaded++
				return
			}

			myComics, _, err := c.GetComicsFromID(comicsID)
			if err != nil && comicsID != 404 {
				log.Printf("%s, comicsID is - %d", err, comicsID)
				return
			}
			if comicsID == 404 {
				successDownloaded++
				return
			}

			mt.Lock()
			defer mt.Unlock()
			c.comics[comicsID] = myComics
			resMap[comicsID] = myComics
			successDownloaded++
		}(i)

	}
	wg.Wait()
	c.lastDownloadedID = c.lastDownloadedID + numOfComics
	return resMap, successDownloaded, nil
}

// TODO: descript
func (c *ComicsDownloader) GetComicsFromID(comicsID int) (core.ComicsDescript, int, error) {
	client := http.Client{}
	comicsURL := fmt.Sprintf("%s/%d/info.0.json", c.comicsURL, comicsID)
	resp, err := client.Get(comicsURL)
	if err != nil {
		return core.ComicsDescript{}, comicsID, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	log.Println(comicsURL)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return core.ComicsDescript{}, comicsID, err
	}

	var fulComics comicsInfo
	err = json.Unmarshal(body, &fulComics)
	if err != nil {
		return core.ComicsDescript{}, comicsID, err
	}

	keywords := strings.Split(fulComics.Transcript, " ")
	return core.ComicsDescript{Url: fulComics.ImgURL, Keywords: keywords}, fulComics.Num, nil
}
