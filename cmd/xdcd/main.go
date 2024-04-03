package main

import (
	"courses/pkg/database"
	"courses/pkg/xkcd"
	"flag"
)

func main() {
	var numOfComics int
	flag.IntVar(&numOfComics, "n", 3, "number of comics to save")
	flag.Parse()
	// TODO: return error
	bytes := xkcd.GetNComicsFromSite("xkcd.com", "db.json", numOfComics)
	database.WriteToDB("db.json", bytes)
}
