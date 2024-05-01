package handler

import (
	"courses/internal/core"
	"courses/internal/core/find"
	"courses/pkg/words"
	"encoding/json"
	"net/http"
	"strings"
)

func CreateServeMux(finder find.Finder) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /pics", func(wr http.ResponseWriter, r *http.Request) {
		comicsKeywords := r.URL.Query().Get("search")
		clearedKeywords := words.StemStringWithClearing(strings.Split(comicsKeywords, " "))
		res := finder.ByIndex(clearedKeywords)
		comicsToSend := min(len(res), core.MaxComicsToShow)
		data, _ := json.Marshal(res[:comicsToSend])
		_, _ = wr.Write(data)
	})
	mux.HandleFunc("POST /update", func(wr http.ResponseWriter, r *http.Request) {
		diff := finder.UpdateComics()
		data, err := json.Marshal(diff)
		if err != nil {
			// TODO
		}
		wr.Write(data)
	})
	return mux
}
