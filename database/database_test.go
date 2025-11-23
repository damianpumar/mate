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

		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		newRecord := records[0].(map[string]interface{})

		if newRecord["id"] != inserted.Id || newRecord["name"] != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, newRecord)
		}
	})

	t.Run("should insert a new value with typed select", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		records := database.SelectTyped[Example](db, "fake")

		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		newRecord := records[0]

		if newRecord.Id != inserted.Id || newRecord.Name != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, newRecord)
		}
	})

	t.Run("should select by id", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		record, found := db.SelectByID("fake", "1")
		if !found {
			t.Fatal("Expected to find record")
		}

		recordMap := record.(map[string]interface{})

		if recordMap["id"] != inserted.Id || recordMap["name"] != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, recordMap)
		}
	})

	t.Run("should select by id with typed select", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		record, found := database.SelectByIDTyped[Example](db, "fake", "1")
		if !found {
			t.Fatal("Expected to find record")
		}

		if record.Id != inserted.Id || record.Name != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, record)
		}
	})

	t.Run("should return false when record not found", func(t *testing.T) {
		defer os.RemoveAll("database")

		_, found := db.SelectByID("fake", "nonexistent")
		if found {
			t.Error("Expected record not to be found")
		}
	})

	t.Run("should update a value", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		updated := Example{Id: "1", Name: "Jane Doe"}

		if !db.Update("fake", inserted.Id, updated) {
			t.Fatal("Failed to update record")
		}

		records := db.Select("fake")

		updatedRecord := records[0].(map[string]interface{})

		if updatedRecord["id"] != updated.Id || updatedRecord["name"] != updated.Name {
			t.Errorf("Expected record %+v, got %+v", updated, updatedRecord)
		}
	})

	t.Run("should update where condition matches", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []Example{
			{Id: "1", Name: "John"},
			{Id: "2", Name: "Jane"},
			{Id: "3", Name: "Bob"},
		}

		for _, ex := range examples {
			db.Insert("fake", ex)
		}

		updated := database.UpdateWhereTyped(db, "fake",
			func(e Example) bool {
				return e.Name[0] == 'J'
			},
			func(e Example) Example {
				e.Name = e.Name + " Updated"
				return e
			})

		if updated != 2 {
			t.Errorf("Expected 2 records updated, got %d", updated)
		}

		records := database.SelectTyped[Example](db, "fake")

		for _, r := range records {
			if r.Name == "John" || r.Name == "Jane" {
				t.Errorf("Expected name to be updated, got %s", r.Name)
			}
		}
	})

	t.Run("should delete a value", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		if !db.Delete("fake", inserted.Id) {
			t.Fatal("Failed to delete record")
		}

		records := db.Select("fake")

		if len(records) != 0 {
			t.Errorf("Expected 0 records, got %d", len(records))
		}
	})

	t.Run("should delete where condition matches", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []Example{
			{Id: "1", Name: "John"},
			{Id: "2", Name: "Jane"},
			{Id: "3", Name: "Bob"},
		}

		for _, ex := range examples {
			db.Insert("fake", ex)
		}

		deleted := database.DeleteWhereTyped(db, "fake", func(e Example) bool {
			return e.Name[0] == 'J'
		})

		if deleted != 2 {
			t.Errorf("Expected 2 records deleted, got %d", deleted)
		}

		records := database.SelectTyped[Example](db, "fake")

		if len(records) != 1 {
			t.Errorf("Expected 1 record remaining, got %d", len(records))
		}

		if records[0].Name != "Bob" {
			t.Errorf("Expected remaining record to be Bob, got %s", records[0].Name)
		}
	})

	t.Run("should count records", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []Example{
			{Id: "1", Name: "John"},
			{Id: "2", Name: "Jane"},
		}

		for _, ex := range examples {
			db.Insert("fake", ex)
		}

		count := db.Count("fake")

		if count != 2 {
			t.Errorf("Expected 2 records, got %d", count)
		}
	})

	t.Run("should check if record exists", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		db.Insert("fake", inserted)

		exists := db.Exists("fake", "1")

		if !exists {
			t.Error("Expected record to exist")
		}

		exists = db.Exists("fake", "999")

		if exists {
			t.Error("Expected record to not exist")
		}
	})

	t.Run("should truncate table", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []Example{
			{Id: "1", Name: "John"},
			{Id: "2", Name: "Jane"},
		}

		for _, ex := range examples {
			db.Insert("fake", ex)
		}

		db.Truncate("fake")

		count := db.Count("fake")
		if count != 0 {
			t.Errorf("Expected 0 records after truncate, got %d", count)
		}
	})

	t.Run("should insert many records", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []interface{}{
			Example{Id: "1", Name: "John"},
			Example{Id: "2", Name: "Jane"},
			Example{Id: "3", Name: "Bob"},
		}

		db.InsertMany("fake", examples)

		count := db.Count("fake")
		if count != 3 {
			t.Errorf("Expected 3 records, got %d", count)
		}
	})

	t.Run("should select where with predicate", func(t *testing.T) {
		defer os.RemoveAll("database")

		examples := []Example{
			{Id: "1", Name: "John"},
			{Id: "2", Name: "Jane"},
			{Id: "3", Name: "Bob"},
		}

		for _, ex := range examples {
			db.Insert("fake", ex)
		}

		results := database.SelectWhereTyped(db, "fake", func(e Example) bool {
			return e.Name[0] == 'J'
		})

		if len(results) != 2 {
			t.Errorf("Expected 2 records, got %d", len(results))
		}
	})

	t.Run("should handle empty table name gracefully", func(t *testing.T) {
		db.Insert("", Example{Id: "1", Name: "Test"})

		records := db.Select("")
		if len(records) != 0 {
			t.Error("Expected empty result for empty table name")
		}
	})
}
