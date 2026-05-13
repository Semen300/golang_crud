package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedBasketItems = []model.TaskCreationDTO{
	{ItemID: 1, Name: "Item1", ItemPrice: 1000, Amount: 10},
	{ItemID: 2, Name: "Item2", ItemPrice: 2000, Amount: 5},
}

func migrateBasketItems(db *sql.DB) {
	_, createErr := db.Exec(`CREATE TABLE IF NOT EXISTS basket_items(
	customerLogin TEXT,
	itemID SERIAL,
	name TEXT,
	itemPrice INTEGER,
	amount INTEGER,
	PRIMARY KEY (customerLogin, itemID))`)
	if createErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'basket_items': \nError creating table: \n%w", createErr))
	}

	_, migrationErr := db.Exec(`INSERT INTO basket_items
	VALUES ($1, $2, $3, $4, $5),
	($6, $7, $8, $9, $10),
	($11, $12, $13, $14, $15),
	($16, $17, $18, $19, $20)`,
		"customer1", 1, "Item1", 1000, 10,
		"customer1", 2, "Item2", 2000, 5,
		"customer2", 1, "Item1", 1000, 10,
		"customer2", 2, "Item2", 2000, 5)
	if migrationErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'basket_items': \nError adding values: \n%w", migrationErr))
	}
}

func resetBasketItems(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS basket_items")
	migrateBasketItems(db)
}

func TestNewBasketRepository(t *testing.T) {
	basketRepo, creatingErr := repository.NewBasketRepository(testDB)
	if creatingErr != nil {
		t.Fatal(creatingErr)
	}

	pingErr := basketRepo.Conn.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
}

func TestGetBasket(t *testing.T) {
	testLogin := "customer1"
	basketRepo, repoErr := repository.NewBasketRepository(testDB)
	assert.Equal(t, nil, repoErr)

	basket, getErr := basketRepo.GetBasket(testLogin)

	assert.Equal(t, nil, getErr)
	assert.Equal(t, expectedBasketItems, basket)
}

func TestSaveToBasketSave(t *testing.T) {
	testLogin, testDTO := "customer1", model.TaskCreationDTO{ItemID: 3, Name: "Item3", ItemPrice: 3000, Amount: 3}
	basketRepo, repoErr := repository.NewBasketRepository(testDB)
	assert.Equal(t, nil, repoErr)

	saveErr := basketRepo.SaveToBasket(testLogin, testDTO)
	assert.Equal(t, nil, saveErr)
	defer resetBasketItems(basketRepo.Conn)

	basket, getErr := basketRepo.GetBasket(testLogin)
	assert.Equal(t, nil, getErr)
	assert.Equal(t, len(basket), 3)
}

func TestSaveToBasketUpdate(t *testing.T) {
	testLogin, testDTO := "customer1", model.TaskCreationDTO{ItemID: 1, Name: "Item1", ItemPrice: 1000, Amount: 3}
	basketRepo, repoErr := repository.NewBasketRepository(testDB)
	assert.Equal(t, nil, repoErr)

	saveErr := basketRepo.SaveToBasket(testLogin, testDTO)
	assert.Equal(t, nil, saveErr)
	defer resetBasketItems(basketRepo.Conn)

	basket, getErr := basketRepo.GetBasket(testLogin)
	assert.Equal(t, nil, getErr)
	assert.Equal(t, 2, len(basket))
	assert.Equal(t, testDTO, basket[0])
	assert.Equal(t, expectedBasketItems[1], basket[1])
}

func TestDeleteFromBasket(t *testing.T) {
	testLogin, testID := "customer1", 1
	basketRepo, repoErr := repository.NewBasketRepository(testDB)
	assert.Equal(t, nil, repoErr)

	basketRepo.DeleteFromBasket(testLogin, testID)
	defer resetBasketItems(testDB)

	basket, getErr := basketRepo.GetBasket(testLogin)
	assert.Equal(t, nil, getErr)
	assert.Equal(t, 1, len(basket))
	assert.Equal(t, expectedBasketItems[1], basket[0])
}

func TestClearBasket(t *testing.T) {
	testLogin := "customer1"
	basketRepo, repoErr := repository.NewBasketRepository(testDB)
	assert.Equal(t, nil, repoErr)

	clearErr := basketRepo.ClearBasket(testLogin)
	if clearErr != nil {
		t.Fatal(clearErr)
	}
	defer resetBasketItems(testDB)

	basket, _ := basketRepo.GetBasket(testLogin)
	assert.Equal(t, 0, len(basket))
}
