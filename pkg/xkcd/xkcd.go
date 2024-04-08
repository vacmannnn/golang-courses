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
	comicsURL        string
	comics           map[int]ComicsDescript
	lastDownloadedID int
}

// NewComicsDownloader sets link to source site with comics
func NewComicsDownloader(comicsURL string, comics map[int]ComicsDescript) ComicsDownloader {
	return ComicsDownloader{comicsURL: comicsURL, comics: comics}
}

// TODO
// func (c *ComicsDownloader) ChangeStartPointOfDownloading(sp int ) {
// 	c.lastDownloadedID = sp
// }

// GetComicsFromSite gets slice of comics indices to download. In case of zero len slice, it will download all possible
// comics. Function will log any non-critical error.
func (c *ComicsDownloader) GetComicsFromSite(numOfComics int) (map[int]ComicsDescript, int, error) {

	var wg sync.WaitGroup
	var mt sync.Mutex
	var successDownloaded int

	func() {
		for i := c.lastDownloadedID + 1; i <= c.lastDownloadedID+numOfComics; i++ {
			fmt.Println(i)
			wg.Add(1)
			go func(comicsID int) {
				defer wg.Done()
				if c.comics[comicsID].Keywords != nil {
					log.Println(comicsID)
					successDownloaded++
					fmt.Println(successDownloaded)
					return
				}
				comicsURL := fmt.Sprintf("%s/%d/info.0.json", c.comicsURL, comicsID)
				log.Println(comicsURL)

				myComics, err := c.getComicsFromURL(comicsURL)
				if err != nil && comicsID != 404 {
					log.Printf("%s, comicsID is - %d", err, comicsID)
					return
				}
				if comicsID == 404 {
					successDownloaded++
					return
				}
				keywords := strings.Split(myComics.Transcript, " ")
				mt.Lock()
				defer mt.Unlock()
				c.comics[comicsID] = ComicsDescript{Url: myComics.ImgURL, Keywords: keywords}
				successDownloaded++
			}(i)

		}
	}()
	wg.Wait()
	c.lastDownloadedID = c.lastDownloadedID + numOfComics
	return c.comics, successDownloaded, nil
}

func (c *ComicsDownloader) getComicsFromURL(comicsURL string) (comicsInfo, error) {
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

func (c *ComicsDownloader) getLatestComicsID() int {
	comicsURL := fmt.Sprintf("%s/info.0.json", c.comicsURL)
	latestComics, _ := c.getComicsFromURL(comicsURL)
	return latestComics.Num
}
