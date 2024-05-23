package database

import (
	"courses/core"
	"math/rand"
	"os"
	"testing"
)

func TestDataBase_WriteAndRead(t *testing.T) {
	comics := make([]core.ComicsDescript, 500)
	db, err := NewDB("testDB.sql", "migration")
	defer db.Close()
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range comics {
		v.Url = "rickroll.com"
		v.Keywords = []string{"rick", "roll"}
		err = db.Write(v, i+1)
		if err != nil {
			t.Error(err)
		}
	}
	newSet, err := db.Read()
	if err != nil {
		t.Error(err)
	}
	if len(newSet) != len(comics) {
		t.Errorf("comic count mismatch, expected: %d, got: %d", len(comics), len(newSet))
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestDataBase_CreateFile(t *testing.T) {
	fileName := RandStringRunes(10)
	_, err := NewDB(fileName, "migration")
	if err != nil {
		t.Error(err)
	}
	err = os.Remove(fileName)
	if err != nil {
		t.Error(err)
	}
}
