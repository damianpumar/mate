package database

import (
	"log"
	"sync"

	"github.com/damianpumar/mate/database/file"
)

type DB struct {
	mu sync.Mutex
}

func Connect() DB {
	return DB{
		mu: sync.Mutex{},
	}
}

func (db *DB) Select(table string) []interface{} {
	data := file.Fetch()

	records := data.Records(table)

	return records
}

func (db *DB) SelectById(table string, id string) interface{} {
	records := db.Select(table)

	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})

		if !ok {
			log.Fatalf("Failed to assert record as map[string]interface{}")
		}

		if recordMap["id"] == id {
			return record
		}
	}

	return nil
}

func (db *DB) Insert(table string, record interface{}) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()

	records := data.Records(table)

	records = append(records, record)

	data.Commit(table, records)

	return true
}

func (db *DB) Update(table string, id string, record interface{}) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()

	records := data.Records(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})

		if !ok {
			log.Fatalf("Failed to assert record as map[string]interface{}")
		}

		if recordMap["id"] == id {
			records[i] = record

			data.Commit(table, records)

			return true
		}
	}

	return false
}

func (db *DB) Delete(table string, id string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()

	records := data.Records(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})

		if !ok {
			log.Fatalf("Failed to assert record as map[string]interface{}")
		}

		if recordMap["id"] == id {
			records = append(records[:i], records[i+1:]...)

			data.Commit(table, records)

			return true
		}
	}

	return false
}
