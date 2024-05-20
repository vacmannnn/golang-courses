package handler

import (
	"courses/internal/core/catalog"
	"log/slog"
	"net/http"
)

type Server struct {
	ctlg   *catalog.ComicsCatalog // interface
	logger *slog.Logger
	mux    *http.ServeMux
}

func NewMux(ctlg *catalog.ComicsCatalog, logger *slog.Logger) http.Handler {
	myMux := Server{ctlg: ctlg, logger: logger, mux: http.NewServeMux()}

	myMux.mux.HandleFunc("GET /pics", myMux.protectedGet())
	myMux.mux.HandleFunc("POST /update", myMux.protectedUpdate())
	myMux.mux.HandleFunc("GET /login", myMux.LoginHandler)

	return limit(myMux.mux)
}
