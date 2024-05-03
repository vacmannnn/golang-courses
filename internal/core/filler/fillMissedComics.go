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
	id int
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
			if err := f.writeComicsWithID(curComics); err != nil {
				return f.comics, err
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
		if f.comics[comID].Keywords == nil {
			descript, id, err := f.downloader.GetComicsFromID(comID)
			if err != nil {
				f.logger.Debug(err.Error(), "comics ID", id)
				results <- comicsDescriptWithID{id: id}
				continue
			}
			f.logger.Info("writing comics", "id", id)
			descript.Keywords = words.StemStringWithClearing(descript.Keywords)
			results <- comicsDescriptWithID{id: id, ComicsDescript: descript}
			mt.Lock()
			f.comics[id] = descript
			mt.Unlock()
			continue
		}
		f.logger.Info("writing comics", "id", comID)
		results <- comicsDescriptWithID{id: comID, ComicsDescript: f.comics[comID]}
	}
}

func (f *Filler) writeComicsWithID(comicsWID comicsDescriptWithID) error {
	var comics = make(map[int]core.ComicsDescript)
	comics[comicsWID.id] = comicsWID.ComicsDescript
	return f.db.Write(comics)
}
