package database

import (
	"log"
	"os"
)

func WriteToDB(dbFileName string, info []byte) {
	// fmt.Println(bytes)
	err := os.WriteFile(dbFileName, info, 0644)
	if err != nil {
		log.Fatal(err)
	}

}
