package repository

import (
	"crud-go/internal/model"
	"database/sql"
	"fmt"
)

type ItemRepository struct {
	Conn      *sql.DB
	CurrentID uint
}

func NewItemRepository(db *sql.DB) (ItemRepository, error) {
	query := `CREATE TABLE IF NOT EXISTS items(
	id SERIAL PRIMARY KEY,
	name TEXT,
	price SERIAL)`

	_, migrationErr := db.Exec(query)
	if migrationErr != nil {
		return ItemRepository{},
			fmt.Errorf("Error creating table \"orders\":\n %w", migrationErr)
	}
	return ItemRepository{db, 0}, nil
}

func (i ItemRepository) GetAllItems() ([]model.Item, error) {
	query := `SELECT *
	FROM items`

	rows, queryErr := i.Conn.Query(query)
	if queryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}
	defer rows.Close()

	var items []model.Item = make([]model.Item, 0)
	for rows.Next() {
		var item model.Item
		scanErr := rows.Scan(&item.Id, &item.Name, &item.Price)
		if scanErr != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanErr)
		}
		items = append(items, item)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}
	return items, nil
}

func (i ItemRepository) GetItemById(id int) (model.Item, error) {
	query := `SELECT *
	FROM items
	WHERE id = $1`

	var item model.Item
	scanErr := i.Conn.QueryRow(query, id).Scan(&item.Id, &item.Name, &item.Price)
	if scanErr != nil {
		return model.Item{},
			fmt.Errorf("Error scanning values from rows:\n %w", scanErr)
	}

	return item, nil
}
