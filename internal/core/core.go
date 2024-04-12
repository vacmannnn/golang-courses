package core

import "time"

type ComicsDescript struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

const MaxWaitTime = time.Second * 3
