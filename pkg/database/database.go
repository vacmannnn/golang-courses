package database

import (
	"os"
)

func WriteToDB(dbFileName string, info []byte) error {
	return os.WriteFile(dbFileName, info, 0644)
}

func ReadFromDB(dbFileName string) ([]byte, error) {
	return os.ReadFile(dbFileName)
}
