package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ComicsData struct {
	ID        int
	URL       string
	ImageName string
}

type SearchMessage struct {
	Message string
}

// todo: handle nil cookie incognito mode http://localhost:3000/comics?search=%27apple+doctor%27
func comicsFinder(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwtTokenCookie")
	if err != nil || cookie == nil {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	var htmlData SearchMessage
	comicsKeywords := r.URL.Query().Get("search")
	if comicsKeywords != "" {
		comicsKeywords = url.QueryEscape(comicsKeywords)
		data, err := sendSearchRequest(comicsKeywords, cookie.Value)
		if err != nil {
			htmlData.Message = "Unexpected problem. Try to search again"
		}
		if len(data) > 0 {
			tmpl, _ := template.ParseFiles("templates/comics_results.html")
			_ = tmpl.Execute(w, data)
			return
		}
		htmlData.Message = "Comics not found"
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		searchString := r.FormValue("comicsToSearch")
		if searchString != "" {
			searchString = url.QueryEscape(searchString)
			http.Redirect(w, r, fmt.Sprintf("/comics?search='%s'", searchString), http.StatusMovedPermanently)
		}
		htmlData.Message = "Enter non-empty string"
	}
	tmpl, _ := template.ParseFiles("templates/comics_search.html")
	_ = tmpl.Execute(w, htmlData)
}

func sendSearchRequest(searchString string, token string) ([]ComicsData, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://localhost:8080/pics?search='%s'", searchString), nil)
	if err != nil {
		log.Printf("creating request: %v\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("sending request: %v\n", err)
		return nil, err
	}
	str, _ := io.ReadAll(res.Body)

	searchResult := make([]string, 0, 10)
	err = json.Unmarshal(str, &searchResult)
	if err != nil {
		log.Printf("unmarshalling result: %v\n", err)
		return nil, err
	}

	var data []ComicsData
	for i, v := range searchResult {
		// name example - https://imgs.xkcd.com/comics/magnet_fishing.png
		data = append(data, ComicsData{ID: i + 1, URL: v, ImageName: v[29 : len(v)-4]})
	}
	return data, nil
}
