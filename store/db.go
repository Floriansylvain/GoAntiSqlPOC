package store

import (
	"log"
	"sync"

	"github.com/dgraph-io/badger/v4"
)

var (
	db     *badger.DB
	once   sync.Once
	dbPath = "./data"
)

func InitDB() {
	once.Do(func() {
		opts := badger.DefaultOptions(dbPath).WithLogger(nil)
		var err error
		db, err = badger.Open(opts)
		if err != nil {
			log.Fatal(err)
		}
	})
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func GetDB() *badger.DB {
	return db
}
