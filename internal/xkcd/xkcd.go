package xkcd

import (
	"courses/internal/core"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	comicsURL string
}

// NewComicsDownloader sets link to source site with comics
func NewComicsDownloader(comicsURL string) ComicsDownloader {
	return ComicsDownloader{comicsURL: comicsURL}
}

// GetComicsFromID downloads comics with given ID. Returns its picture url, transcript and id.
func (c *ComicsDownloader) GetComicsFromID(comicsID int) (core.ComicsDescript, int, error) {
	if comicsID == 404 {
		return core.ComicsDescript{Url: "https://xkcd.com/404", Keywords: nil}, comicsID, nil
	}

	client := http.Client{Timeout: core.MaxWaitTime}
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

	keywords := strings.Split(fulComics.Transcript+fulComics.Alt, " ")
	return core.ComicsDescript{Url: fulComics.ImgURL, Keywords: keywords}, fulComics.Num, nil
}
