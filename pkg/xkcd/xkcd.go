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
// описать случай, при котором error != nil, но при этом в byte был записан json, но там может быть кривой json
func GetNComicsFromSite(urlName string, dbFileName string, comicsNum int) ([]byte, error) {
	if comicsNum < 1 {
		return nil, errors.New("number of comics should be greater than 0, default value is 1")
	}
	comicsToJSON := make(map[int]comicsDescript)
	lastNum := 1
	file, err := database.ReadFromDB(dbFileName)
	// если ошибка, то ничего страшного, т.к. все равно перезапишем весь файл потом, а сейчас чтение для попытки
	// не делать скачку дважды
	if err != nil {
		log.Println(err)
	}

	// теоретически, если файл был кривой, то ничего страшного, перезапишем все
	// TODO: найти случай, при котором вообще ошибка появляется
	err = json.Unmarshal(file, &comicsToJSON)
	if err != nil {
		log.Println(err)
		lastNum = 1
	} else {
		for k := range comicsToJSON {
			lastNum = max(lastNum, k)
		}
	}

	// last comicsInfo will be overwritten due possible corruption
	for i := lastNum; i <= comicsNum; i++ {
		c := http.Client{}
		comicsURL := fmt.Sprintf("https://%s/%d/info.0.json", urlName, i)
		log.Println(comicsURL)
		resp, err := c.Get(comicsURL)
		if err != nil {
			bytes, _ := marshallComics(comicsToJSON)
			return bytes, err
		}
		// defer in loop !! maybe close explicitly ?
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			bytes, _ := marshallComics(comicsToJSON)
			return bytes, err
		}
		var myComics comicsInfo
		err = json.Unmarshal(body, &myComics)
		if err != nil {
			return nil, err
		}
		keywords := words.StemStringWithClearing(myComics.Transcript)
		comicsToJSON[myComics.Num] = comicsDescript{Url: myComics.ImgURL, Keywords: keywords}
	}

	bytes, err := marshallComics(comicsToJSON)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func marshallComics(comicsToJSON map[int]comicsDescript) ([]byte, error) {
	bytes, err := json.MarshalIndent(comicsToJSON, "", " ")
	return bytes, err
}
