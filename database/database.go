package database

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/damianpumar/mate/database/file"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidRecord  = errors.New("invalid record format")
	ErrEmptyTable     = errors.New("table name cannot be empty")
)

type DB struct {
	mu sync.RWMutex
}

func Connect() *DB {
	return &DB{
		mu: sync.RWMutex{},
	}
}

func convertToType[T any](record interface{}) (T, error) {
	var result T

	if typedRecord, ok := record.(T); ok {
		return typedRecord, nil
	}

	if recordMap, ok := record.(map[string]interface{}); ok {
		jsonBytes, err := json.Marshal(recordMap)
		if err != nil {
			return result, err
		}
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return result, err
		}
		return result, nil
	}

	return result, ErrInvalidRecord
}

func (db *DB) Select(table string) []interface{} {
	if table == "" {
		return []interface{}{}
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	return data.Records(table)
}

func SelectTyped[T any](db *DB, table string) []T {
	records := db.Select(table)
	result := make([]T, 0, len(records))

	for _, record := range records {
		if converted, err := convertToType[T](record); err == nil {
			result = append(result, converted)
		}
	}

	return result
}

func (db *DB) SelectByID(table string, id string) (interface{}, bool) {
	if table == "" {
		return nil, false
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	records := data.Records(table)

	for _, record := range records {
		recordMap, ok := record.(map[string]interface{})
		if !ok {
			continue
		}

		if recordMap["id"] == id {
			return record, true
		}
	}

	return nil, false
}

func SelectByIDTyped[T any](db *DB, table string, id string) (*T, bool) {
	record, found := db.SelectByID(table, id)
	if !found {
		return nil, false
	}

	converted, err := convertToType[T](record)
	if err != nil {
		return nil, false
	}

	return &converted, true
}

func (db *DB) SelectWhere(table string, predicate func(interface{}) bool) []interface{} {
	records := db.Select(table)
	result := make([]interface{}, 0)

	for _, record := range records {
		if predicate(record) {
			result = append(result, record)
		}
	}

	return result
}

func SelectWhereTyped[T any](db *DB, table string, predicate func(T) bool) []T {
	records := SelectTyped[T](db, table)
	result := make([]T, 0)

	for _, record := range records {
		if predicate(record) {
			result = append(result, record)
		}
	}

	return result
}

func (db *DB) Insert(table string, record interface{}) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)
	records = append(records, record)

	data.Commit(table, records)
}

func (db *DB) InsertMany(table string, records []interface{}) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	existingRecords := data.Records(table)
	existingRecords = append(existingRecords, records...)

	data.Commit(table, existingRecords)
}

func InsertManyTyped[T any](db *DB, table string, records []T) {
	interfaceRecords := make([]interface{}, len(records))
	for i, record := range records {
		interfaceRecords[i] = record
	}
	db.InsertMany(table, interfaceRecords)
}

func (db *DB) Upsert(table string, id string, record interface{}) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		if recordMap["id"] == id {
			records[i] = record
			data.Commit(table, records)
			return
		}
	}

	records = append(records, record)
	data.Commit(table, records)
}

func UpsertTyped[T any](db *DB, table string, id string, record T) {
	db.Upsert(table, id, record)
}

func (db *DB) Update(table string, id string, record interface{}) bool {
	if table == "" {
		return false
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		if recordMap["id"] == id {
			records[i] = record
			data.Commit(table, records)
			return true
		}
	}

	return false
}

func (db *DB) UpsertMany(table string, records []interface{}) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	existingRecords := data.Records(table)

	recordMap := make(map[string]int)
	for i, r := range existingRecords {
		if rMap, ok := r.(map[string]interface{}); ok {
			if id, ok := rMap["id"].(string); ok {
				recordMap[id] = i
			}
		}
	}

	for _, newRecord := range records {
		newRecordMap, ok := newRecord.(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := newRecordMap["id"].(string)
		if !ok {
			continue
		}

		if idx, exists := recordMap[id]; exists {
			existingRecords[idx] = newRecord
		} else {
			existingRecords = append(existingRecords, newRecord)
			recordMap[id] = len(existingRecords) - 1
		}
	}

	data.Commit(table, existingRecords)
}

func UpsertManyTyped[T any](db *DB, table string, records []T) {
	interfaceRecords := make([]interface{}, len(records))
	for i, record := range records {
		interfaceRecords[i] = record
	}
	db.UpsertMany(table, interfaceRecords)
}

func (db *DB) UpsertWhere(table string, predicate func(interface{}) bool, record interface{}, updater func(interface{}) interface{}) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		if predicate(r) {
			records[i] = updater(r)
			data.Commit(table, records)
			return
		}
	}

	records = append(records, record)
	data.Commit(table, records)
}

func UpsertWhereTyped[T any](db *DB, table string, predicate func(T) bool, record T, updater func(T) T) {
	db.UpsertWhere(table,
		func(r interface{}) bool {
			if converted, err := convertToType[T](r); err == nil {
				return predicate(converted)
			}
			return false
		},
		record,
		func(r interface{}) interface{} {
			if converted, err := convertToType[T](r); err == nil {
				return updater(converted)
			}
			return r
		})
}

func (db *DB) UpdateWhere(table string, predicate func(interface{}) bool, updater func(interface{}) interface{}) int {
	if table == "" {
		return 0
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)
	updated := 0

	for i, r := range records {
		if predicate(r) {
			records[i] = updater(r)
			updated++
		}
	}

	if updated > 0 {
		data.Commit(table, records)
	}

	return updated
}

func UpdateWhereTyped[T any](db *DB, table string, predicate func(T) bool, updater func(T) T) int {
	return db.UpdateWhere(table,
		func(record interface{}) bool {
			if converted, err := convertToType[T](record); err == nil {
				return predicate(converted)
			}
			return false
		},
		func(record interface{}) interface{} {
			if converted, err := convertToType[T](record); err == nil {
				return updater(converted)
			}
			return record
		})
}

func (db *DB) Delete(table string, id string) bool {
	if table == "" {
		return false
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		if recordMap["id"] == id {
			records = append(records[:i], records[i+1:]...)
			data.Commit(table, records)
			return true
		}
	}

	return false
}

func (db *DB) DeleteWhere(table string, predicate func(interface{}) bool) int {
	if table == "" {
		return 0
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)
	deleted := 0

	newRecords := make([]interface{}, 0, len(records))
	for _, r := range records {
		if !predicate(r) {
			newRecords = append(newRecords, r)
		} else {
			deleted++
		}
	}

	if deleted > 0 {
		data.Commit(table, newRecords)
	}

	return deleted
}

func DeleteWhereTyped[T any](db *DB, table string, predicate func(T) bool) int {
	return db.DeleteWhere(table, func(record interface{}) bool {
		if converted, err := convertToType[T](record); err == nil {
			return predicate(converted)
		}
		return false
	})
}

func (db *DB) Count(table string) int {
	if table == "" {
		return 0
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	return len(data.Records(table))
}

func (db *DB) Exists(table string, id string) bool {
	if table == "" {
		return false
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	records := data.Records(table)

	for _, r := range records {
		recordMap, ok := r.(map[string]interface{})
		if !ok {
			continue
		}

		if recordMap["id"] == id {
			return true
		}
	}

	return false
}

func (db *DB) Truncate(table string) {
	if table == "" {
		return
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	data.Commit(table, []interface{}{})
}

func (db *DB) Transaction(fn func(*DB) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return fn(db)
}

func (db *DB) Map(table string, mapper func(interface{}) interface{}) []interface{} {
	records := db.Select(table)
	result := make([]interface{}, 0, len(records))

	for _, record := range records {
		result = append(result, mapper(record))
	}

	return result
}

func MapTyped[T any, R any](db *DB, table string, mapper func(T) R) []R {
	records := SelectTyped[T](db, table)
	result := make([]R, 0, len(records))

	for _, record := range records {
		result = append(result, mapper(record))
	}

	return result
}
