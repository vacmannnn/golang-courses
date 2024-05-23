package database

import (
	"courses/core"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

type DataBase struct {
	pathToDB string
	db       *sql.DB
}

// NewDB sets path to database
func NewDB(path string, pathToMigrations string) (*DataBase, error) {
	err := createFileIfNotExists(path)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	err = runMigrate(db, pathToMigrations)
	if err != nil {
		return nil, err
	}
	return &DataBase{pathToDB: path, db: db}, nil
}

func runMigrate(db *sql.DB, migrationsDir string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite3 db driver failed %s", err)
	}

	dir := fmt.Sprintf("file://%s", migrationsDir)
	fileSource, err := (&file.File{}).Open(dir)
	if err != nil {
		return fmt.Errorf("opening migration file failed %s", err)
	}

	m, err := migrate.NewWithInstance("file", fileSource, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("initializing db migration failed %s", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrating database failed %s", err)
	}

	return nil
}

func createFileIfNotExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, create it.
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				_ = file.Close()
			}(file)
		} else {
			return err
		}
	}
	return nil
}

func (d *DataBase) Write(descript core.ComicsDescript, comicsID int) error {
	keywords := strings.Join(descript.Keywords, " ")
	_, err := d.db.Exec("insert or ignore into comics (url, keywords, comicsID) values ($1, $2, $3)",
		descript.Url, keywords, comicsID)

	return err
}

func (d *DataBase) Read() (map[int]core.ComicsDescript, error) {
	rows, err := d.db.Query("select * from comics")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	comics := make(map[int]core.ComicsDescript)

	for rows.Next() {
		var id int
		var keywords string
		descript := core.ComicsDescript{}
		err = rows.Scan(&descript.Url, &keywords, &id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		descript.Keywords = strings.Split(keywords, " ")

		comics[id] = descript
	}

	return comics, err
}

func (d *DataBase) Close() {
	_ = d.db.Close()
}
