package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
)

type MyData struct {
	ID  int
	URL string
}

func comicsFinder(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwtTokenCookie")
	if err != nil {
		// TODO: redirect to login
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	}

	if cookie != nil {
		fmt.Println(cookie.Value)
	}

	comicsKeywords := r.URL.Query().Get("search")
	fmt.Println(comicsKeywords)
	if comicsKeywords != "" {
		comicsKeywords = url.QueryEscape(comicsKeywords)
		fmt.Println("here", r.Method, comicsKeywords)
		req, err := http.NewRequest(http.MethodGet,
			fmt.Sprintf("http://localhost:8080/pics?search='%s'", comicsKeywords), nil)
		if err != nil {
			log.Printf("creating request: %v\n", err)
		}
		req.Header.Set("Authorization", cookie.Value)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		str, _ := io.ReadAll(res.Body)

		searchResult := make([]string, 0, 10)
		err = json.Unmarshal(str, &searchResult)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(searchResult)

		var data []MyData

		for i, v := range searchResult {
			data = append(data, MyData{ID: i + 1, URL: v})
		}
		fmt.Println(data)
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, data)
		return
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		searchString := r.FormValue("comicsToSearch")
		log.Println(searchString)
		searchString = url.QueryEscape(searchString)
		log.Println(searchString)
		http.Redirect(w, r, fmt.Sprintf("/comics?search='%s'", searchString), http.StatusMovedPermanently)
	} else {
		http.ServeFile(w, r, "templates/comics.html")
	}
}
