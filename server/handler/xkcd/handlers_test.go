package xkcd

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestLogin(t *testing.T) {
	// set handler and launch server which will be shut down by context
	hndlr := NewServerHandler(&catalog{}, *logger, "../../users.json", 10, 10, 10)
	s := httptest.NewServer(hndlr)
	defer s.Close()

	// send request with admins pass and name
	var jsonStr = []byte(`{"username":"admin", "password":"admin"}`)
	req, err := http.NewRequest(http.MethodPost, s.URL+"/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status: got %v", res.Status)
	}
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
}

func TestSearchNotFoundRequest(t *testing.T) {
	// set handler and launch server which will be shut down by context
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil,
		tokenMaxTime: 0, requests: make(chan struct{})}
	srv.mux.HandleFunc("GET /pics", func(w http.ResponseWriter, r *http.Request) {
		srv.searchRequest(w, r)
	})
	s := httptest.NewServer(srv.mux)
	defer s.Close()

	cases := []struct {
		name           string
		searchString   string
		expectedStatus int
	}{
		{
			name:           "comics won't found",
			searchString:   "/pics?search='banana'",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "comics will be found",
			searchString:   "/pics?search='abc'",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "empty search string",
			searchString:   "/pics",
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// send search request
			res, err := http.Get(s.URL + tc.searchString)
			if err != nil {
				t.Fatalf("sending request: %v", err)
			}

			// handle response
			if res.StatusCode != tc.expectedStatus {
				t.Errorf("unexpected status: got %v", res.Status)
			}
		})
	}
}

func TestUpdateRequest(t *testing.T) {
	// set handler and launch server which will be shut down by context
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil, tokenMaxTime: 1, requests: make(chan struct{})}
	srv.mux.HandleFunc("POST /update", func(w http.ResponseWriter, r *http.Request) {
		srv.updateRequest(w, r)
	})
	s := httptest.NewServer(srv.mux)
	defer s.Close()

	// send search request to update db
	res, err := http.Post(s.URL+"/update", "", nil)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	// handle response
	_, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cannot read from response body: %v", err)
	}
}

func TestAuth(t *testing.T) {
	srv := server{ctlg: &catalog{}, logger: *logger, mux: http.NewServeMux(), users: nil, tokenMaxTime: 10, requests: make(chan struct{}, 1)}
	srv.mux.HandleFunc("GET /pics", srv.protectedUpdate())
	s := httptest.NewServer(srv.mux)
	defer s.Close()

	cases := []struct {
		name           string
		user           userInfo
		needToken      bool
		token          string
		expectedStatus int
	}{
		{name: "valid data",
			user: userInfo{Username: "user",
				Password: "$2a$10$w6/HvzjDEJa7vgmEGWtXCuz9YkUkcyLMHN547wRhNyUTR0zPIILmK"},
			needToken:      true,
			expectedStatus: http.StatusOK,
		}, {
			name: "no token",
			user: userInfo{Username: "user",
				Password: "$2a$10$w6/HvzjDEJa7vgmEGWtXCuz9YkUkcyLMHN547wRhNyUTR0zPIILmK"},
			needToken:      false,
			expectedStatus: http.StatusUnauthorized,
		}, {
			name: "incorrect token",
			user: userInfo{Username: "user",
				Password: "$2a$10$w6/HvzjDEJa7vgmEGWtXCuz9YkUkcyLMHN547wRhNyUTR0zPIILmK"},
			needToken:      true,
			token:          "abc.abc.abc",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.needToken && tc.token == "" {
				tc.token, _ = createToken(userInfo{Username: "user",
					Password: "$2a$10$w6/HvzjDEJa7vgmEGWtXCuz9YkUkcyLMHN547wRhNyUTR0zPIILmK"}, 10)
			}
			req, err := http.NewRequest(http.MethodGet, s.URL+"/pics?search='abc'", nil)
			if err != nil {
				t.Errorf("creating request: %v", err)
			}
			if tc.needToken {
				req.Header.Set("Authorization", tc.token)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf("sending request: %v", err)
			}
			if res.StatusCode != tc.expectedStatus {
				t.Errorf("unexpected response status: %v", res.StatusCode)
			}
		})
	}
}
