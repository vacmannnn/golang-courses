package handler

import (
	"courses/core"
	"courses/pkg/words"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

// TODO: handle errors on fprint

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	// concurrency limiter
	s.requests <- struct{}{}
	defer func() { <-s.requests }()

	w.Header().Set("Content-Type", "application/json")

	var u userInfo
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		// TODO
		return
	}

	role, err := s.auth(u)

	if err == nil {
		u.role = role
		tokenString, err := createToken(u, s.tokenMaxTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Error("No Username found")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid credentials")
	}
}

func (s *server) auth(user userInfo) (int, error) {
	for _, u := range s.users {
		if u.Username == user.Username &&
			bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)) == nil {
			return u.role, nil
		}
	}
	return -1, fmt.Errorf("invalid user")
}

func (s *server) protectHandler(next func(http.ResponseWriter, *http.Request), checkForAdmin bool,
	w http.ResponseWriter, r *http.Request) {
	// concurrency limiter
	s.requests <- struct{}{}
	defer func() { <-s.requests }()

	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}

	isAdmin, err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	if checkForAdmin && !isAdmin {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Only admin can update comics set")
		return
	}

	next(w, r)
}

func (s *server) searchRequest(wr http.ResponseWriter, r *http.Request) {
	comicsKeywords := r.URL.Query().Get("search")
	if comicsKeywords == "" {
		wr.WriteHeader(http.StatusNotFound)
		_, err := wr.Write([]byte("404 empty search string"))
		if err != nil {
			s.logger.Error("writing response for GET /pics", "err", err)
		}
		return
	}

	clearedKeywords := words.StemStringWithClearing(strings.Split(comicsKeywords, " "))
	res := s.ctlg.FindByIndex(clearedKeywords)

	if len(res) == 0 {
		wr.WriteHeader(404)
		_, err := wr.Write([]byte("404 not found"))
		if err != nil {
			s.logger.Error("writing response for GET /pics", "err", err)
		}
		return
	}

	comicsToSend := min(len(res), core.MaxComicsToShow)
	data, err := json.Marshal(res[:comicsToSend])
	if err != nil {
		s.logger.Error("marshalling res of catalog", "err", err)
		data = []byte("")
	}

	_, err = wr.Write(data)
	if err != nil {
		s.logger.Error("writing response for GET /pics", "err", err)
	}
}

func (s *server) protectedSearch() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.protectHandler(s.searchRequest, false, w, r)
	}
}

func (s *server) updateRequest(wr http.ResponseWriter, _ *http.Request) {
	numOfNewComics, total, err := s.ctlg.UpdateComics()
	if err != nil {
		s.logger.Error("updating comics", "err", err)
	}

	diff := map[string]int{"new": numOfNewComics, "total": total}
	data, err := json.Marshal(diff)
	if err != nil {
		s.logger.Error("marshalling diff of comics update", "err", err)
		data = []byte("")
	}

	_, err = wr.Write(data)
	if err != nil {
		s.logger.Error("writing response for POST /update", "err", err)
	}
}

func (s *server) protectedUpdate() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s.protectHandler(s.updateRequest, true, w, r)
	}
}
