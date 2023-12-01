package main

import (
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/timshannon/badgerhold/v4"
)

type MetaData struct {
	Record map[string]any
}

type File struct {
	Name string
}

type ProvenanceData struct {
	Dataset string
	Created time.Time
	Files   []File
}

type BadgerRecord interface {
	MetaData | ProvenanceData
}

func StoreInit(dir string) (*badgerhold.Store, error) {
	options := badgerhold.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	store, err := badgerhold.Open(options)
	log.Printf("store %+v, error %v\n", store, err)
	return store, err
}

// Insert inserts badger records records
func Insert[T BadgerRecord](store *badgerhold.Store, records []T) error {
	err := store.Badger().Update(func(tx *badger.Txn) error {
		for _, rec := range records {
			uuid, _ := uuid.NewRandom()
			err := store.TxInsert(tx, uuid, rec)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Find finds records in our store for given query
func Find[T BadgerRecord](store *badgerhold.Store, query *badgerhold.Query) ([]T, error) {
	var records []T
	err := store.Find(&records, query)
	return records, err
}
