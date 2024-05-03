package handler

import (
	"courses/internal/core"
	"courses/internal/core/catalog"
	"courses/pkg/words"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func CreateServeMux(ctlg *catalog.ComicsCatalog, logger *slog.Logger) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /pics", func(wr http.ResponseWriter, r *http.Request) {
		comicsKeywords := r.URL.Query().Get("search")
		clearedKeywords := words.StemStringWithClearing(strings.Split(comicsKeywords, " "))
		res := ctlg.FindByIndex(clearedKeywords)

		if len(res) == 0 {
			wr.WriteHeader(404)
			_, err := wr.Write([]byte("404 not found"))
			if err != nil {
				logger.Error("writing response for GET /pics", "err", err)
			}
			return
		}

		comicsToSend := min(len(res), core.MaxComicsToShow)
		data, err := json.Marshal(res[:comicsToSend])
		if err != nil {
			logger.Error("marshalling res of catalog", "err", err)
			data = []byte("")
		}
		_, err = wr.Write(data)
		if err != nil {
			logger.Error("writing response for GET /pics", "err", err)
		}
	})

	mux.HandleFunc("POST /update", func(wr http.ResponseWriter, r *http.Request) {
		diff, err := ctlg.UpdateComics()
		if err != nil {
			logger.Error("updating comics", "err", err)
		}
		data, err := json.Marshal(diff)
		if err != nil {
			logger.Error("marshalling diff of comics update", "err", err)
			data = []byte("")
		}
		_, err = wr.Write(data)
		if err != nil {
			logger.Error("writing response for POST /update", "err", err)
		}
	})

	return mux
}
