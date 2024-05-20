package handler

import (
	"courses/internal/core"
	"log/slog"
	"net/http"
)

type Server struct {
	ctlg   core.Catalog
	logger slog.Logger
	mux    *http.ServeMux
}

func NewMux(ctlg core.Catalog, logger slog.Logger) http.Handler {
	myServ := Server{ctlg: ctlg, logger: logger, mux: http.NewServeMux()}

	myServ.mux.HandleFunc("GET /pics", myServ.protectedSearch())
	myServ.mux.HandleFunc("POST /update", myServ.protectedUpdate())
	myServ.mux.HandleFunc("POST /login", myServ.login)

	return limit(myServ.mux)
}
