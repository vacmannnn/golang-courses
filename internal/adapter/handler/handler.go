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

type Mux struct {
	ctlg         *catalog.ComicsCatalog // interface
	logger       *slog.Logger
	idkHowToName *http.ServeMux
}

func NewMux(ctlg *catalog.ComicsCatalog, logger *slog.Logger) http.Handler {
	myMux := Mux{ctlg: ctlg, logger: logger, idkHowToName: http.NewServeMux()}

	myMux.idkHowToName.HandleFunc("GET /pics", protectedGet(myMux.UpdateRequestHandler))
	myMux.idkHowToName.HandleFunc("POST /update", protectedUpdate(myMux.GetRequestHandler))
	myMux.idkHowToName.HandleFunc("GET /login", LoginHandler)

	return limit(myMux.idkHowToName)
}

func (m *Mux) GetRequestHandler(wr http.ResponseWriter, r *http.Request) {
	comicsKeywords := r.URL.Query().Get("search")
	//validate?
	clearedKeywords := words.StemStringWithClearing(strings.Split(comicsKeywords, " "))
	res := m.ctlg.FindByIndex(clearedKeywords)

	if len(res) == 0 {
		wr.WriteHeader(404)
		_, err := wr.Write([]byte("404 not found"))
		if err != nil {
			m.logger.Error("writing response for GET /pics", "err", err)
		}
		return
	}

	comicsToSend := min(len(res), core.MaxComicsToShow)
	data, err := json.Marshal(res[:comicsToSend])
	if err != nil {
		m.logger.Error("marshalling res of catalog", "err", err)
		data = []byte("")
	}
	_, err = wr.Write(data)
	if err != nil {
		m.logger.Error("writing response for GET /pics", "err", err)
	}
}

func (m *Mux) UpdateRequestHandler(wr http.ResponseWriter, _ *http.Request) {
	diff, err := m.ctlg.UpdateComics()
	if err != nil {
		m.logger.Error("updating comics", "err", err)
	}
	data, err := json.Marshal(diff)
	if err != nil {
		m.logger.Error("marshalling diff of comics update", "err", err)
		data = []byte("")
	}
	_, err = wr.Write(data)
	if err != nil {
		m.logger.Error("writing response for POST /update", "err", err)
	}
}
