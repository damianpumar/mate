package database_test

import (
	"os"
	"testing"

	"github.com/damianpumar/mate/database"
)

func TestDatabase(t *testing.T) {
	type Example struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	db := database.Connect()

	t.Run("should insert a new value", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		records := db.Select("fake")

		newRecord := records[0].(map[string]interface{})

		if newRecord["id"] != inserted.Id && newRecord["name"] != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, newRecord)
		}
	})

	t.Run("should update a value", func(t *testing.T) {
		defer os.RemoveAll("database")
		inserted := Example{Id: "1", Name: "John Doe"}

		if ok := db.Insert("fake", inserted); !ok {
			t.Errorf("Failed to insert record")
		}

		updated := Example{Id: "1", Name: "Jane Doe"}

		if ok := db.Update("fake", inserted.Id, updated); !ok {
			t.Errorf("Failed to update record")
		}

		records := db.Select("fake")

		updatedRecord := records[0].(map[string]interface{})

		if updatedRecord["id"] != updated.Id && updatedRecord["name"] != updated.Name {
			t.Errorf("Expected record %+v, got %+v", updated, updatedRecord)
		}
	})

	t.Run("should delete a value", func(t *testing.T) {
		defer os.RemoveAll("database")
		inserted := Example{Id: "1", Name: "John Doe"}

		if ok := db.Insert("fake", inserted); !ok {
			t.Errorf("Failed to insert record")
		}

		if ok := db.Delete("fake", inserted.Id); !ok {
			t.Errorf("Failed to delete record")
		}

		records := db.Select("fake")

		if len(records) != 0 {
			t.Errorf("Expected 0 records, got %d", len(records))
		}
	})
}
