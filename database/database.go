package database

import (
	"sync"

	"github.com/damianpumar/mate/database/file"
)

type Identifiable interface {
	GetId() string
}

type DB struct {
	mu sync.Mutex
}

func Connect() DB {
	return DB{
		mu: sync.Mutex{},
	}
}

func (db *DB) Select(table string) []interface{} {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)
	return records
}

func SelectById[T Identifiable](db *DB, table string, id string) (T, bool) {
	records := SelectBy[T](db, table, func(record T) bool {
		return record.GetId() == id
	})

	if len(records) > 0 {
		return records[0], true
	}

	var zero T
	return zero, false
}

func SelectBy[T Identifiable](db *DB, table string, predicate func(T) bool) []T {
	db.mu.Lock()
	defer db.mu.Unlock()

	records := db.Select(table)
	var results []T

	for _, record := range records {
		typedRecord, ok := record.(T)
		if !ok {
			continue
		}

		if predicate(typedRecord) {
			results = append(results, typedRecord)
		}
	}

	return results
}

func Insert[T Identifiable](db *DB, table string, record T) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	records = append(records, record)

	data.Commit(table, records)

	return true
}

func Update[T Identifiable](db *DB, table string, id string, record T) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		typedRecord, ok := r.(T)
		if !ok {
			continue
		}

		if typedRecord.GetId() == id {
			records[i] = record
			data.Commit(table, records)
			return true
		}
	}

	return false
}

func Upsert[T Identifiable](db *DB, table string, record T) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		typedRecord, ok := r.(T)
		if !ok {
			continue
		}

		if typedRecord.GetId() == record.GetId() {
			records[i] = record
			data.Commit(table, records)
			return true
		}
	}

	records = append(records, record)
	data.Commit(table, records)

	return true
}

func Drop(db *DB, table string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	data.Commit(table, []interface{}{})
}

func Delete[T Identifiable](db *DB, table string, id string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		typedRecord, ok := r.(T)
		if !ok {
			continue
		}

		if typedRecord.GetId() == id {
			records = append(records[:i], records[i+1:]...)
			data.Commit(table, records)
			return true
		}
	}

	return false
}
