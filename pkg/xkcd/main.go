package main

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

const comicsNum = 10

type comicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

func main() {
	comicsToJSON := make(map[int]comicsDescript, comicsNum)
	for i := 1; i < comicsNum; i++ {
		c := http.Client{}
		comicsURL := fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)
		fmt.Println(comicsURL)
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
	fmt.Println(comicsToJSON)
	bytes, err := json.MarshalIndent(comicsToJSON, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bytes)
	fileName := "db.json"
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		return
	}
}
