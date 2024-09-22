package database

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type DB struct {
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

	var db DB
	err = json.Unmarshal(byteValue, &db)

	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %s", err)
	}

	return db
}

func (db DB) Save() {
	file, err := os.OpenFile("database/database.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}

	defer file.Close()

	byteValue, err := json.MarshalIndent(db, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %s", err)
	}

	_, err = file.Write(byteValue)
	if err != nil {
		log.Fatalf("Failed to write to file: %s", err)
	}
}
