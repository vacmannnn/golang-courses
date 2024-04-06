package database

import (
	"os"
)

type DataBase struct {
	pathToDB string
}

// NewDB sets path to database
func NewDB(path string) DataBase {
	return DataBase{pathToDB: path}
}

func (d DataBase) Write(info []byte) error {
	return os.WriteFile(d.pathToDB, info, 0644)
}

func (d DataBase) Read() ([]byte, error) {
	return os.ReadFile(d.pathToDB)
}
