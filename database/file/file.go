package file

import (
	"encoding/json"
	"log"
	"os"
)

type Data struct {
	data interface{}
}

const PATH = "database/database.json"

func Fetch() *Data {
	if _, err := os.Stat(PATH); os.IsNotExist(err) {
		file, err := os.Create(PATH)

		if err != nil {
			log.Fatalf("Failed to create database: %s", err)
		}

		file.Close()
	}

	byteValue, err := os.ReadFile(PATH)

	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	var data interface{}

	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	return &Data{
		data,
	}
}

func (f *Data) Commit(key string, value interface{}) {
	dataMap, ok := f.data.(map[string]interface{})

	if !ok {
		log.Fatalf("Failed to assert data as map[string]interface{}")
	}

	dataMap[key] = value

	file, err := os.OpenFile(PATH, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	byteValue, err := json.MarshalIndent(f.data, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatalf("Failed to write to file: %s", err)
	}
}

func (f *Data) Records(table string) []interface{} {
	dataMap, ok := f.data.(map[string]interface{})

	if !ok {
		log.Fatalf("Failed to assert data as map[string]interface{}")
	}

	if dataMap[table] == nil {
		return make([]interface{}, 0)
	}

	return dataMap[table].([]interface{})
}
