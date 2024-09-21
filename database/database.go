package database

import (
	"encoding/json"
	"log"

	"github.com/patrickmn/go-cache"
	"github.com/shubhexists/go-json-db/models"
)

type Database struct {
	db    *models.Driver
	cache *cache.Cache
}

func Connect() Database {
	db, cache, err := models.New("./db")

	if err != nil {
		log.Fatal(err)
	}

	return Database{db, cache}
}

func (d *Database) Set(key string, value interface{}) {
	d.db.Write(key, value)
}

func (d *Database) Get(key string) []interface{} {
	records, err := d.db.ReadAll(key, d.cache, false)

	if err != nil {
		log.Fatal(err)
	}

	results := make([]interface{}, len(records))

	for k, v := range records {
		var result interface{}

		err := json.Unmarshal([]byte(v), &result)

		if err != nil {
			log.Fatal(err)
		}

		results[k] = result
	}

	return results
}
