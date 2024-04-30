package handler

import (
	"courses/internal/core"
	"courses/internal/core/find"
	"courses/pkg/words"
	"encoding/json"
	"net/http"
	"reflect"
	"slices"
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
		updatedComics, err := filler.FillMissedComics()
		if err != nil {
			// TODO
		}
		eq := reflect.DeepEqual(updatedComics, comics)
		var data []byte
		var n int
		if !eq {
			for k, v := range updatedComics {
				if slices.Equal(comics[k].Keywords, v.Keywords) {
					n++
				}
			}
			// TODO: shared memory, case with everyday update
			comics = updatedComics
		}
		diff := map[string]int{
			"new": n, "total": len(updatedComics),
		}
		data, err = json.Marshal(diff)
		if err != nil {
			// TODO
		}
		wr.Write(data)
	})
	return mux
}
