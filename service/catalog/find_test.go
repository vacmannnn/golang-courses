package catalog

import (
	"context"
	"courses/core"
	"courses/service/filler"
	"courses/service/xkcd"
	database "courses/storage"
	"io"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"testing"
)

type myFiller struct {
}

func (f *myFiller) FillMissedComics(context.Context) (map[int]core.ComicsDescript, error) {
	curComics := make(map[int]core.ComicsDescript)
	for k, v := range comics {
		curComics[k] = v
	}
	curComics[5] = core.ComicsDescript{Url: "https://imgs.xkcd.com/comics/blownapart_color.jpg", Keywords: []string{"black",
		"red", "packag", "hey", "packag", "packag", "explod", "boom", "red", "cloud",
		"smoke", "red", "green", "blue", "lie", "scorch", "mark", "floor", "blown", "prime", "factor", "blown",
		"prime", "factor"}}
	return curComics, nil
}

var comics = map[int]core.ComicsDescript{
	1: {Url: "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", Keywords: []string{"boy", "sit", "barrel", "float",
		"ocean", "boy", "float", "barrel", "drift", "distanc"}},
	2: {Url: "https://imgs.xkcd.com/comics/tree_cropped_(1).jpg", Keywords: []string{"tree", "grow", "sphere", "titl",
		"petit", "refer", "petit", "princ", "halfway", "sketch", "petit", "refer", "petit", "princ",
		"halfway", "sketch"}},
	3: {Url: "https://imgs.xkcd.com/comics/island_color.jpg", Keywords: []string{"sketch", "petit", "island", "island"}},
	4: {Url: "https://imgs.xkcd.com/comics/landscape_cropped_(1).jpg", Keywords: []string{"sketch", "landscap", "horizon",
		"river", "flow", "ocean", "river", "flow", "ocean"}},
}

var cases = []struct {
	name              string
	searchString      []string
	expectedToBeFound bool
}{
	{
		name:              "should be found",
		searchString:      strings.Split("tree grow", " "),
		expectedToBeFound: true,
	},
	{
		name:              "shouldn't be found",
		searchString:      strings.Split("red dress", " "),
		expectedToBeFound: false,
	},
	{
		name:              "empty string",
		searchString:      []string{},
		expectedToBeFound: false,
	},
}

func TestSearchByIndex(t *testing.T) {
	ctlg := NewCatalog(comics, &myFiller{})
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := ctlg.FindByIndex(tc.searchString)
			if len(res) > 0 != tc.expectedToBeFound {
				if tc.expectedToBeFound {
					t.Errorf("expected to be found, but found %d comics", len(res))
				} else {
					t.Errorf("unexpected to be found, but found %d comics", len(res))
				}
			}
		})
	}
}

func TestSearchByComics(t *testing.T) {
	ctlg := NewCatalog(comics, &myFiller{})
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := ctlg.findByComics(tc.searchString)
			if len(res) > 0 != tc.expectedToBeFound {
				if tc.expectedToBeFound {
					t.Errorf("expected to be found, but found %d comics", len(res))
				} else {
					t.Errorf("unexpected to be found, but found %d comics", len(res))
				}
			}
		})
	}
}

func TestUpdateComics(t *testing.T) {
	ctlg := NewCatalog(comics, &myFiller{})

	newComics, total, _ := ctlg.UpdateComics()
	if newComics != 1 || total != 5 {
		t.Errorf("should be founded new comics, but num of new comics: %v, total: %v", newComics, total)
	}
}

func BenchmarkDiffMethToSearch(b *testing.B) {
	myDB, _ := database.NewDB("test.json", "/storage/migration")

	comics, _ := myDB.Read()
	if comics == nil {
		comics = make(map[int]core.ComicsDescript, 3000)
	}

	sourceUrl := "xkcd.com"
	downloader := xkcd.NewComicsDownloader(sourceUrl)

	opts := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(io.Discard, opts)
	logger := slog.New(handler)
	comicsFiller := filler.NewFiller(core.GoroutineNum, comics, myDB, downloader, *logger)
	comics, _ = comicsFiller.FillMissedComics(context.Background())

	index := make(map[string][]int)
	var doc []string
	for k, v := range comics {
		doc = slices.Concat(doc, v.Keywords)
		for i, token := range v.Keywords {
			if !slices.Contains(v.Keywords[:i], token) {
				index[token] = append(index[token], k)
			}
		}
	}
	testString := []string{"my favorite comics is about unknown mystery person", "idk what comics to search",
		"cool banana man", "orange box sits under that orange table and takes orange to make orange juice",
		"funny comics about math"}
	for _, str := range testString {
		comicsName := "findFindByIndex-" + strconv.Itoa(len(str))
		catalog := NewCatalog(comics, &comicsFiller)
		b.Run(comicsName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				catalog.FindByIndex(strings.Split(str, " "))
			}
		})
		comicsName = "findByComics-" + strconv.Itoa(len(str))
		b.Run(comicsName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				catalog.findByComics(strings.Split(str, " "))
			}
		})
	}
}
