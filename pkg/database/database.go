package database

import (
	"os"
)

// TODO: description
func WriteToDB(dbFileName string, info []byte) error {
	return os.WriteFile(dbFileName, info, 0644)
}

func ReadFromDB(dbFileName string) ([]byte, error) {
	// TODO: what if empty file or file doesnt exist ?
	return os.ReadFile(dbFileName)
}
