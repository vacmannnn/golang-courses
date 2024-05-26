package xkcd

import (
	"errors"
	"testing"
)

func TestComicsDownloader_GetComicsFromID(t *testing.T) {
	downloader := NewComicsDownloader("https://xkcd.com")

	testCases := []struct {
		name     string
		comicsID int
		err      error
	}{
		{
			name:     "404 comics",
			comicsID: 404,
			err:      nil,
		},
		{
			name:     "existing comics",
			comicsID: 100,
			err:      nil,
		},
		{
			name:     "not existing comics",
			comicsID: 10001,
			err:      errors.New("some error"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, id, err := downloader.GetComicsFromID(tc.comicsID)

			// tc.err and err both should be nil or not nil
			if (err == nil) != (tc.err == nil) {
				t.Errorf("GetComicsFromID(%d): %v", tc.comicsID, err)
			}

			if id != tc.comicsID {
				t.Errorf("Unexpected comics ID: %d != %d", id, tc.comicsID)
			}
		})
	}
}
