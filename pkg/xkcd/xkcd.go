package xkcd

import (
	"courses/pkg/words"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

type comicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// TODO: description
func GetNComicsFromSite(urlName string, dbFileName string, comicsNum int) []byte {
	// TODO: handle case if comicsNum < 1
	comicsToJSON := make(map[int]comicsDescript)
	lastNum := 1
	file, _ := os.ReadFile(dbFileName)
	json.Unmarshal(file, &comicsToJSON)
	for k := range comicsToJSON {
		lastNum = max(lastNum, k)
	}
	// fmt.Println(comicsToJSON)
	// log.Fatal("abc")
	// last comicsInfo will be overwritten due possible corruption
	// fmt.Println(lastNum, comicsNum)
	for i := lastNum; i <= comicsNum; i++ {
		c := http.Client{}
		comicsURL := fmt.Sprintf("https://%s/%d/info.0.json", urlName, i)
		log.Println(comicsURL)
		resp, err := c.Get(comicsURL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var myComics comicsInfo
		err = json.Unmarshal(body, &myComics)
		if err != nil {
			log.Fatal(err)
		}
		keywords := words.StemStringWithClearing(myComics.Transcript)
		comicsToJSON[myComics.Num] = comicsDescript{Url: myComics.ImgURL, Keywords: keywords}
	}
	// fmt.Println(comicsToJSON)
	bytes, err := json.MarshalIndent(comicsToJSON, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
