package filler

import (
	"context"
	"courses/internal/core"
	"courses/pkg/words"
	"log/slog"
	"sync"
)

type Filler struct {
	goroutineNum int
	comics       map[int]core.ComicsDescript
	db           core.DataBase
	downloader   core.ComicsDownloader
	logger       slog.Logger
}

type comicsDescriptWithID struct {
	core.ComicsDescript
	id           int
	isDownloaded bool
}

func NewFiller(goroutineNum int, comics map[int]core.ComicsDescript, db core.DataBase,
	downloader core.ComicsDownloader, logger slog.Logger) Filler {
	return Filler{
		goroutineNum: goroutineNum,
		comics:       comics,
		db:           db,
		downloader:   downloader,
		logger:       logger,
	}
}

func (f *Filler) FillMissedComics(ctx context.Context) (map[int]core.ComicsDescript, error) {
	comicsIDChan := make(chan int, f.goroutineNum)
	comicsChan := make(chan comicsDescriptWithID, f.goroutineNum)
	wg := sync.WaitGroup{}
	var mt sync.Mutex

	done := ctx.Done()
	closed := func() bool {
		select {
		case <-done:
			return true
		default:
			return false
		}
	}

	// launch worker pool
	for range f.goroutineNum {
		wg.Add(1)
		go func() {
			f.worker(comicsIDChan, comicsChan, &mt)
			wg.Done()
		}()
	}

	var curComics comicsDescriptWithID
	// download comics till no error
	for i := 1; ; i++ {
		// send in advance bunch of ID to optimize downloading
		if i%f.goroutineNum == 1 {
			for j := i; j < i+f.goroutineNum; j++ {
				comicsIDChan <- j
			}
		}

		curComics = <-comicsChan
		if curComics.Url != "" && !closed() {
			if curComics.isDownloaded {
				f.logger.Info("writing comics to DB", "id", curComics.id)
				if err := f.writeComicsWithID(curComics); err != nil {
					return f.comics, err
				}
			}
		} else {
			close(comicsIDChan)
			wg.Wait()
			close(comicsChan)

			for range len(comicsChan) {
				curComics = <-comicsChan
				if curComics.Url == "" || !curComics.isDownloaded {
					continue
				}
				f.logger.Info("writing comics to DB", "id", curComics.id)
				if err := f.writeComicsWithID(curComics); err != nil {
					return f.comics, err
				}
			}
			break
		}
	}
	return f.comics, nil
}

func (f *Filler) worker(comicsIDChan <-chan int, results chan<- comicsDescriptWithID, mt *sync.Mutex) {
	for comID := range comicsIDChan {
		f.logger.Info("working on comics", "id", comID)
		if f.comics[comID].Keywords == nil {
			descript, id, err := f.downloader.GetComicsFromID(comID)
			if err != nil {
				f.logger.Debug(err.Error(), "comics ID", id)
				results <- comicsDescriptWithID{id: id}
				continue
			}
			descript.Keywords = words.StemStringWithClearing(descript.Keywords)
			results <- comicsDescriptWithID{id: id, ComicsDescript: descript, isDownloaded: true}
			mt.Lock()
			f.comics[id] = descript
			mt.Unlock()
			continue
		}
		results <- comicsDescriptWithID{id: comID, ComicsDescript: f.comics[comID]}
	}
}

func (f *Filler) writeComicsWithID(comicsWID comicsDescriptWithID) error {
	return f.db.Write(comicsWID.ComicsDescript, comicsWID.id)
}
