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
	myServ := Server{ctlg: ctlg, logger: logger, mux: http.NewServeMux()}

	myServ.mux.HandleFunc("GET /pics", myServ.protectedSearch())
	myServ.mux.HandleFunc("POST /update", myServ.protectedUpdate())
	myServ.mux.HandleFunc("POST /login", myServ.login)

	return limit(myServ.mux)
}
