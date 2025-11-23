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

func (db *DB) Select(table string) ([]interface{}, error) {
	if table == "" {
		return nil, ErrEmptyTable
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	records := data.Records(table)

	return records, nil
}

func SelectTyped[T any](db *DB, table string) ([]T, error) {
	records, err := db.Select(table)
	if err != nil {
		return nil, err
	}

	result := make([]T, 0, len(records))
	for _, record := range records {
		if converted, err := convertToType[T](record); err == nil {
			result = append(result, converted)
		}
	}

	return result, nil
}

func (db *DB) SelectByID(table string, id string) (interface{}, error) {
	if table == "" {
		return nil, ErrEmptyTable
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
			return record, nil
		}
	}

	return nil, ErrRecordNotFound
}

func SelectByIDTyped[T any](db *DB, table string, id string) (T, error) {
	var zero T

	record, err := db.SelectByID(table, id)
	if err != nil {
		return zero, err
	}

	return convertToType[T](record)
}

func (db *DB) SelectWhere(table string, predicate func(interface{}) bool) ([]interface{}, error) {
	records, err := db.Select(table)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0)
	for _, record := range records {
		if predicate(record) {
			result = append(result, record)
		}
	}

	return result, nil
}

func SelectWhereTyped[T any](db *DB, table string, predicate func(T) bool) ([]T, error) {
	records, err := SelectTyped[T](db, table)
	if err != nil {
		return nil, err
	}

	result := make([]T, 0)
	for _, record := range records {
		if predicate(record) {
			result = append(result, record)
		}
	}

	return result, nil
}

func (db *DB) Insert(table string, record interface{}) error {
	if table == "" {
		return ErrEmptyTable
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)
	records = append(records, record)

	data.Commit(table, records)
	return nil
}

func (db *DB) InsertMany(table string, records []interface{}) error {
	if table == "" {
		return ErrEmptyTable
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	existingRecords := data.Records(table)

	existingRecords = append(existingRecords, records...)

	data.Commit(table, existingRecords)
	return nil
}

func InsertManyTyped[T any](db *DB, table string, records []T) error {
	interfaceRecords := make([]interface{}, len(records))
	for i, record := range records {
		interfaceRecords[i] = record
	}
	return db.InsertMany(table, interfaceRecords)
}

func (db *DB) Upsert(table string, id string, record interface{}) error {
	if table == "" {
		return ErrEmptyTable
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
			return nil
		}
	}

	records = append(records, record)
	data.Commit(table, records)
	return nil
}

func UpsertTyped[T any](db *DB, table string, id string, record T) error {
	return db.Upsert(table, id, record)
}

func (db *DB) Update(table string, id string, record interface{}) error {
	if table == "" {
		return ErrEmptyTable
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
			return nil
		}
	}

	return ErrRecordNotFound
}

func (db *DB) UpsertMany(table string, records []interface{}) error {
	if table == "" {
		return ErrEmptyTable
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	existingRecords := data.Records(table)

	// Crear un mapa para búsqueda rápida de registros existentes
	recordMap := make(map[string]int)
	for i, r := range existingRecords {
		if rMap, ok := r.(map[string]interface{}); ok {
			if id, ok := rMap["id"].(string); ok {
				recordMap[id] = i
			}
		}
	}

	// Procesar cada registro nuevo
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
			// Actualizar registro existente
			existingRecords[idx] = newRecord
		} else {
			// Agregar nuevo registro
			existingRecords = append(existingRecords, newRecord)
			recordMap[id] = len(existingRecords) - 1
		}
	}

	data.Commit(table, existingRecords)
	return nil
}

func UpsertManyTyped[T any](db *DB, table string, records []T) error {
	interfaceRecords := make([]interface{}, len(records))
	for i, record := range records {
		interfaceRecords[i] = record
	}
	return db.UpsertMany(table, interfaceRecords)
}

func (db *DB) UpsertWhere(table string, predicate func(interface{}) bool, record interface{}, updater func(interface{}) interface{}) error {
	if table == "" {
		return ErrEmptyTable
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	records := data.Records(table)

	for i, r := range records {
		if predicate(r) {
			// Actualizar registro existente usando el updater
			records[i] = updater(r)
			data.Commit(table, records)
			return nil
		}
	}

	records = append(records, record)
	data.Commit(table, records)
	return nil
}

func UpsertWhereTyped[T any](db *DB, table string, predicate func(T) bool, record T, updater func(T) T) error {
	return db.UpsertWhere(table,
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

func (db *DB) UpdateWhere(table string, predicate func(interface{}) bool, updater func(interface{}) interface{}) (int, error) {
	if table == "" {
		return 0, ErrEmptyTable
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

	return updated, nil
}

func UpdateWhereTyped[T any](db *DB, table string, predicate func(T) bool, updater func(T) T) (int, error) {
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

func (db *DB) Delete(table string, id string) error {
	if table == "" {
		return ErrEmptyTable
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
			return nil
		}
	}

	return ErrRecordNotFound
}

func (db *DB) DeleteWhere(table string, predicate func(interface{}) bool) (int, error) {
	if table == "" {
		return 0, ErrEmptyTable
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

	return deleted, nil
}

func DeleteWhereTyped[T any](db *DB, table string, predicate func(T) bool) (int, error) {
	return db.DeleteWhere(table, func(record interface{}) bool {
		if converted, err := convertToType[T](record); err == nil {
			return predicate(converted)
		}
		return false
	})
}

func (db *DB) Count(table string) (int, error) {
	if table == "" {
		return 0, ErrEmptyTable
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	data := file.Fetch()
	records := data.Records(table)

	return len(records), nil
}

func (db *DB) Exists(table string, id string) (bool, error) {
	if table == "" {
		return false, ErrEmptyTable
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
			return true, nil
		}
	}

	return false, nil
}

func (db *DB) Truncate(table string) error {
	if table == "" {
		return ErrEmptyTable
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	data := file.Fetch()
	data.Commit(table, []interface{}{})

	return nil
}

func (db *DB) Transaction(fn func(*DB) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return fn(db)
}

func (db *DB) Map(table string, mapper func(interface{}) interface{}) ([]interface{}, error) {
	records, err := db.Select(table)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, mapper(record))
	}

	return result, nil
}

func MapTyped[T any, R any](db *DB, table string, mapper func(T) R) ([]R, error) {
	records, err := SelectTyped[T](db, table)
	if err != nil {
		return nil, err
	}

	result := make([]R, 0, len(records))
	for _, record := range records {
		result = append(result, mapper(record))
	}

	return result, nil
}
