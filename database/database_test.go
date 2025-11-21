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

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		records, err := db.Select("fake")
		if err != nil {
			t.Fatalf("Failed to select records: %v", err)
		}

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

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		records, err := database.SelectTyped[Example](db, "fake")
		if err != nil {
			t.Fatalf("Failed to select records: %v", err)
		}

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

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		record, err := db.SelectByID("fake", "1")
		if err != nil {
			t.Fatalf("Failed to select record by ID: %v", err)
		}

		recordMap := record.(map[string]interface{})

		if recordMap["id"] != inserted.Id || recordMap["name"] != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, recordMap)
		}
	})

	t.Run("should select by id with typed select", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		record, err := database.SelectByIDTyped[Example](db, "fake", "1")
		if err != nil {
			t.Fatalf("Failed to select record by ID: %v", err)
		}

		if record.Id != inserted.Id || record.Name != inserted.Name {
			t.Errorf("Expected record %+v, got %+v", inserted, record)
		}
	})

	t.Run("should return error when record not found", func(t *testing.T) {
		defer os.RemoveAll("database")

		_, err := db.SelectByID("fake", "nonexistent")
		if err != database.ErrRecordNotFound {
			t.Errorf("Expected ErrRecordNotFound, got %v", err)
		}
	})

	t.Run("should update a value", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		updated := Example{Id: "1", Name: "Jane Doe"}

		if err := db.Update("fake", inserted.Id, updated); err != nil {
			t.Fatalf("Failed to update record: %v", err)
		}

		records, err := db.Select("fake")
		if err != nil {
			t.Fatalf("Failed to select records: %v", err)
		}

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
			if err := db.Insert("fake", ex); err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
		}

		updated, err := database.UpdateWhereTyped(db, "fake",
			func(e Example) bool {
				return e.Name[0] == 'J'
			},
			func(e Example) Example {
				e.Name = e.Name + " Updated"
				return e
			})

		if err != nil {
			t.Fatalf("Failed to update records: %v", err)
		}

		if updated != 2 {
			t.Errorf("Expected 2 records updated, got %d", updated)
		}

		records, _ := database.SelectTyped[Example](db, "fake")

		for _, r := range records {
			if r.Name == "John" || r.Name == "Jane" {
				t.Errorf("Expected name to be updated, got %s", r.Name)
			}
		}
	})

	t.Run("should delete a value", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		if err := db.Delete("fake", inserted.Id); err != nil {
			t.Fatalf("Failed to delete record: %v", err)
		}

		records, err := db.Select("fake")
		if err != nil {
			t.Fatalf("Failed to select records: %v", err)
		}

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
			if err := db.Insert("fake", ex); err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
		}

		deleted, err := database.DeleteWhereTyped(db, "fake", func(e Example) bool {
			return e.Name[0] == 'J'
		})

		if err != nil {
			t.Fatalf("Failed to delete records: %v", err)
		}

		if deleted != 2 {
			t.Errorf("Expected 2 records deleted, got %d", deleted)
		}

		records, _ := database.SelectTyped[Example](db, "fake")

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
			if err := db.Insert("fake", ex); err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
		}

		count, err := db.Count("fake")
		if err != nil {
			t.Fatalf("Failed to count records: %v", err)
		}

		if count != 2 {
			t.Errorf("Expected 2 records, got %d", count)
		}
	})

	t.Run("should check if record exists", func(t *testing.T) {
		defer os.RemoveAll("database")

		inserted := Example{Id: "1", Name: "John Doe"}

		if err := db.Insert("fake", inserted); err != nil {
			t.Fatalf("Failed to insert record: %v", err)
		}

		exists, err := db.Exists("fake", "1")
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

		if !exists {
			t.Error("Expected record to exist")
		}

		exists, err = db.Exists("fake", "999")
		if err != nil {
			t.Fatalf("Failed to check existence: %v", err)
		}

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
			if err := db.Insert("fake", ex); err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
		}

		if err := db.Truncate("fake"); err != nil {
			t.Fatalf("Failed to truncate table: %v", err)
		}

		count, _ := db.Count("fake")
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

		if err := db.InsertMany("fake", examples); err != nil {
			t.Fatalf("Failed to insert many records: %v", err)
		}

		count, _ := db.Count("fake")
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
			if err := db.Insert("fake", ex); err != nil {
				t.Fatalf("Failed to insert record: %v", err)
			}
		}

		results, err := database.SelectWhereTyped(db, "fake", func(e Example) bool {
			return e.Name[0] == 'J'
		})

		if err != nil {
			t.Fatalf("Failed to select where: %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 records, got %d", len(results))
		}
	})

	t.Run("should return error for empty table name", func(t *testing.T) {
		err := db.Insert("", Example{Id: "1", Name: "Test"})
		if err != database.ErrEmptyTable {
			t.Errorf("Expected ErrEmptyTable, got %v", err)
		}
	})
}
