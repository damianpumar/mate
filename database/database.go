package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

type DB struct {
	mu   sync.Mutex
	data interface{}
}

func Connect() DB {
	file, err := os.Open("database/database.json")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(io.Reader(file))

	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var data interface{}
	err = json.Unmarshal(byteValue, &data)

	if err != nil {
		log.Fatalf("Failed to load data: %s", err)
	}

	return DB{
		data: data,
		mu:   sync.Mutex{},
	}
}

func (db *DB) Commit() {
	file, err := os.OpenFile("database/database.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	byteValue, err := json.MarshalIndent(db.data, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatalf("Failed to write to file: %s", err)
	}
}

func (db *DB) Select(table string) []interface{} {
	dataMap, ok := db.data.(map[string]interface{})

	if !ok {
		log.Fatalf("Failed to assert data as map[string]interface{}")
	}

	if dataMap[table] == nil {
		return make([]interface{}, 0)
	}

	return dataMap[table].([]interface{})
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

func (db *DB) Insert(table string, record interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	records := db.Select(table)

	if len(records) == 0 {
		records = make([]interface{}, 0)
	}

	records = append(records, record)

	dataMap, ok := db.data.(map[string]interface{})

	if !ok {
		log.Fatalf("Failed to assert data as map[string]interface{}")
	}

	dataMap[table] = records

	db.Commit()
}

func (db *DB) Update(table string, id string, record interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()

	records := db.Select(table)

	for i, r := range records {
		recordMap, ok := r.(map[string]interface{})

		if !ok {
			log.Fatalf("Failed to assert record as map[string]interface{}")
		}

		if recordMap["id"] == id {
			records[i] = record
			break
		}
	}

	dataMap, ok := db.data.(map[string]interface{})

	if !ok {
		log.Fatalf("Failed to assert data as map[string]interface{}")
	}

	dataMap[table] = records

	db.Commit()
}
