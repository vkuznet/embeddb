package main

import (
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/timshannon/badgerhold/v4"
)

type Item struct {
	ID       int
	Category string `badgerholdIndex:"Category"`
	Created  time.Time
}

type MetaRecord struct {
	ID   int
	Meta map[string]any
}

func Example() {
	data := []Item{
		{
			ID:       0,
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		{
			ID:       1,
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		{
			ID:       2,
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:       3,
			Category: "blue",
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}

	mrec := make(map[string]any)
	mrec["bla"] = "foo"
	mrec["int"] = 1
	metaData := []MetaRecord{
		{
			ID:   1,
			Meta: mrec,
		},
		{
			ID:   2,
			Meta: mrec,
		},
	}

	dir := tempdir()
	log.Println("tempdir", dir)
	defer os.RemoveAll(dir)

	options := badgerhold.DefaultOptions
	options.Dir = dir
	options.ValueDir = dir
	store, err := badgerhold.Open(options)
	log.Printf("store %+v\n", store)
	defer store.Close()

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// insert the data in one transaction

	err = store.Badger().Update(func(tx *badger.Txn) error {
		for i := range data {
			log.Printf("insert record %+v", data)
			err := store.TxInsert(tx, data[i].ID, data[i])
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	// Find all items in the blue category that have been created in the past hour
	var result []Item

	err = store.Find(&result, badgerhold.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1*time.Hour)))

	if err != nil {
		// handle error
		log.Fatal(err)
	}
	log.Printf("looked up %d records for Category=blue\n", len(result))
	for _, rec := range result {
		log.Printf("record %+v", rec)
	}

	//     fmt.Println(result[0].ID)

	// look-up all records
	var records []Item
	query := &badgerhold.Query{}

	err = store.Find(&records, query)

	if err != nil {
		// handle error
		log.Fatal(err)
	}
	log.Printf("looked up %d records for no condition\n", len(records))
	for _, rec := range records {
		log.Printf("record %+v", rec)
	}

	// insert metadata
	err = store.Badger().Update(func(tx *badger.Txn) error {
		for i := range metaData {
			log.Printf("insert meta-record %+v", metaData)
			err := store.TxInsert(tx, metaData[i].ID, metaData[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	// look-up all records
	var metaRecords []MetaRecord
	err = store.Find(&metaRecords, query)

	if err != nil {
		// handle error
		log.Fatal(err)
	}
	log.Printf("looked up %d records for no condition\n", len(metaRecords))
	for _, rec := range metaRecords {
		log.Printf("record %+v", rec)
	}

}
