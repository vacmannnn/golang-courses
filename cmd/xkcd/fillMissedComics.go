package main

import (
	"courses/internal/core"
	"courses/internal/database"
	"courses/internal/xkcd"
	"courses/pkg/words"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func fillMissedComics(goroutineNum int, comics map[int]core.ComicsDescript,
	db database.DataBase, downloader xkcd.ComicsDownloader) (map[int]core.ComicsDescript, error) {

	comicsIDChan := make(chan int, goroutineNum)
	comicsChan := make(chan comicsDescriptWithID, goroutineNum)
	wg := sync.WaitGroup{}
	var mt sync.Mutex

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// launch worker pool
	for range goroutineNum {
		wg.Add(1)
		go func() {
			worker(downloader, comics, comicsIDChan, comicsChan, &mt)
			wg.Done()
		}()
	}

	var curComics comicsDescriptWithID
	// download comics till no error
	for i := 1; ; i++ {
		// send in advance bunch of ID to optimize downloading
		if i%goroutineNum == 1 {
			for j := i; j < i+goroutineNum; j++ {
				comicsIDChan <- j
			}
		}

		curComics = <-comicsChan
		if curComics.Url != "" && len(sigs) == 0 {
			if err := writeComicsWithID(curComics, &db); err != nil {
				return comics, err
			}
		} else {
			close(comicsIDChan)
			wg.Wait()
			close(comicsChan)

			for range len(comicsChan) {
				curComics = <-comicsChan
				if curComics.Url == "" {
					continue
				}
				if err := writeComicsWithID(curComics, &db); err != nil {
					return comics, err
				}
			}
			break
		}
	}
	return comics, nil
}

func worker(downloader xkcd.ComicsDownloader, comics map[int]core.ComicsDescript, comicsIDChan <-chan int,
	results chan<- comicsDescriptWithID, mt *sync.Mutex) {
	for comID := range comicsIDChan {
		if comics[comID].Keywords == nil {
			descript, id, err := downloader.GetComicsFromID(comID)
			if err != nil {
				log.Println(err)
				results <- comicsDescriptWithID{id: comID}
				continue
			}
			descript.Keywords = words.StemStringWithClearing(descript.Keywords)
			results <- comicsDescriptWithID{id: id, ComicsDescript: descript}
			mt.Lock()
			comics[id] = descript
			mt.Unlock()
			continue
		}
		results <- comicsDescriptWithID{id: comID, ComicsDescript: comics[comID]}
	}
}

func writeComicsWithID(comicsWID comicsDescriptWithID, db *database.DataBase) error {
	var comics = make(map[int]core.ComicsDescript)
	comics[comicsWID.id] = comicsWID.ComicsDescript
	return db.Write(comics)
}
