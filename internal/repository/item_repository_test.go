package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

var expectedItems []model.Item = []model.Item{
	{Id: 1, Name: "item1", Price: 100},
	{Id: 2, Name: "item2", Price: 200},
	{Id: 3, Name: "item3", Price: 300},
}

func migrateItems(db *sql.DB) {
	_, createErr := db.Exec(`CREATE TABLE IF NOT EXISTS items(
	id SERIAL PRIMARY KEY,
	name TEXT,
	price SERIAL)`)

	if createErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'items': \nError creating table: \n%w", createErr))
	}

	_, queryErr := db.Exec(`INSERT INTO items VALUES
	($1, $2, $3),
	($4, $5, $6),
	($7, $8, $9)`,
		1, "item1", 100,
		2, "item2", 200,
		3, "item3", 300)
	if queryErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'items': \nError adding values: \n%w", createErr))
	}
}

// func resetItems(db *sql.DB) {
// 	db.Exec(`DROP TABLE IF EXISTS items`)
// 	migrateItems(db)
// }

func TestNewItemRepository(t *testing.T) {
	itemRepo, creatingErr := repository.NewItemRepository(testDB)
	if creatingErr != nil {
		t.Fatal(creatingErr)
	}

	pingErr := itemRepo.Conn.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
}

func TestGetAllItems(t *testing.T) {
	itemRepo, repoErr := repository.NewItemRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}

	items, getErr := itemRepo.GetAllItems()
	if getErr != nil {
		t.Fatal(getErr)
	}

	if len(items) != len(expectedItems) {
		t.Errorf("Result length missmatch: expected %d, got %d", len(expectedItems), len(items))
	}

	for i, item := range items {
		if item != expectedItems[i] {
			t.Errorf("Result missmatch: expected %s, got %s", expectedItems[i].ToString(), item.ToString())
		}
	}
}

func TestGetItemById(t *testing.T) {
	var rowNum = 1
	itemRepo, repoErr := repository.NewItemRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}

	item, getErr := itemRepo.GetItemById(rowNum)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if item != expectedItems[rowNum-1] {
		t.Errorf("Result missmatch: expected %s, got %s", expectedItems[rowNum-1].ToString(), item.ToString())
	}
}
