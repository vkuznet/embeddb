package main

// Copyright 2019 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/timshannon/badgerhold/v4"
)

func tempdir() string {
	name, err := ioutil.TempDir("", "badgerhold-")
	if err != nil {
		panic(err)
	}
	return name
}

func Badger() {
	dir := tempdir()
	log.Println("tempdir", dir)
	defer os.RemoveAll(dir)

	store, err := StoreInit(dir)
	defer store.Close()

	rec1 := make(map[string]any)
	rec1["bla"] = "foo"
	rec1["int"] = 1
	rec2 := make(map[string]any)
	rec2["foo"] = "bla"
	rec2["int"] = 0
	metaData := []MetaData{
		{Record: rec1},
		{Record: rec2},
	}

	// insert meta-data records
	err = Insert(store, metaData)
	if err != nil {
		log.Fatal(err)
	}
	query := &badgerhold.Query{}
	records, err := Find[MetaData](store, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("query %+v found %d records, error %v\n", query, len(records), err)
	for _, rec := range records {
		log.Printf("record %+v", rec)
	}

	// insert Provenance records
	files := []File{
		{Name: "file1.root"},
		{Name: "file2.root"},
	}
	data := []ProvenanceData{
		{
			Dataset: "/a/b/c",
			Created: time.Now().Add(-4 * time.Hour),
			Files:   files,
		},
		{
			Dataset: "/x/y/z",
			Created: time.Now(),
			Files:   files,
		},
	}
	err = Insert(store, data)
	if err != nil {
		log.Fatal(err)
	}
	precords, err := Find[ProvenanceData](store, query)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("query %+v found %d records, error %v\n", query, len(precords), err)
	for _, rec := range precords {
		log.Printf("record %+v", rec)
	}
}

func main() {
	//     Example()
	Badger()
}
