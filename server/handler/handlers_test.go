package handler

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

type catalog struct{}

func (c *catalog) GetIndex() map[string][]int {
	return map[string][]int{}
}

func (c *catalog) UpdateComics() (int, int, error) {
	return 100, 2000, nil
}

func (c *catalog) FindByIndex(input []string) []string {
	if len(input) > 0 && input[0] == "abc" {
		return []string{"abc.com"}
	}
	return nil
}

var logger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func TestMain(m *testing.M) {
	hndlr := NewMux(&catalog{}, *logger, "../../users.json", 10, 10, 10)
	go func() {
		err := http.ListenAndServe("localhost:8080", hndlr)
		if err != nil {
			return
		}
	}()
	// wait for server to launch
	time.Sleep(time.Second * 2)
	m.Run()
}

func TestLogin(t *testing.T) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	var jsonStr = []byte(`{"username":"admin", "password":"admin"}`)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: got %v", res.Status)
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
	// Print the body
	fmt.Println(string(token))
}

func TestSearchNotFoundRequest(t *testing.T) {
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil, tokenMaxTime: 0, requests: make(chan struct{})}
	srv.mux.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		srv.searchRequest(w, r)
	})
	go func() { http.ListenAndServe("localhost:8070", srv.mux) }()
	time.Sleep(time.Second * 3)

	res, err := http.Get("http://localhost:8070/pics?search='banana'")
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status: got %v", res.Status)
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
	// Print the body
	fmt.Println(string(token))
}

func TestSearchRequest(t *testing.T) {
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil, tokenMaxTime: 0, requests: make(chan struct{})}
	srv.mux.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		srv.searchRequest(w, r)
	})
	go func() { http.ListenAndServe("localhost:8050", srv.mux) }()
	time.Sleep(time.Second * 3)

	res, err := http.Get("http://localhost:8050/pics?search='abc'")
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: got %v", res.Status)
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
	// Print the body
	fmt.Println(string(token))
}

func TestUpdateRequest(t *testing.T) {
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil, tokenMaxTime: 1, requests: make(chan struct{})}
	srv.mux.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		srv.updateRequest(w, r)
	})
	go func() { http.ListenAndServe("localhost:8060", srv.mux) }()
	time.Sleep(time.Second * 3)

	res, err := http.Post("http://localhost:8060/update", "", nil)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: got %v", res.Status)
	}
	token, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
	// Print the body
	fmt.Println(string(token))
}
