package database

import (
	"courses/pkg/xkcd"
	"encoding/json"
	"os"
)

type DataBase struct {
	pathToDB   string
	readBefore bool
}

// NewDB sets path to database
func NewDB(path string) DataBase {
	return DataBase{pathToDB: path}
}

func (d *DataBase) Write(data map[int]xkcd.ComicsDescript) error {
	var err error
	w := &os.File{}
	if d.readBefore {
		w, err = os.OpenFile(d.pathToDB, os.O_WRONLY, 0644)
	} else {
		w, err = os.Create(d.pathToDB)
		d.readBefore = true
	}
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func (d *DataBase) Read() (map[int]xkcd.ComicsDescript, error) {
	r, err := os.OpenFile(d.pathToDB, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(r)
	total := make(map[int]xkcd.ComicsDescript)
	for decoder.More() {
		var person map[int]xkcd.ComicsDescript
		if err := decoder.Decode(&person); err != nil {
			return nil, err
		}
		for k, v := range person {
			total[k] = v
		}
	}

	return total, nil
}
