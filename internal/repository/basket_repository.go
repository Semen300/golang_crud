package repository

import (
	"context"
	"crud-go/internal/model"
	"database/sql"
	"fmt"
)

type IBasketRepository interface {
	GetBasket(string) ([]model.TaskCreationDTO, error)
	SaveToBasket(string, model.TaskCreationDTO) error
	DeleteFromBasket(string, int) error

	Begin() (*sql.Tx, error)
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

type BasketRepository struct {
	Conn *sql.DB
}

func NewBasketRepository(db *sql.DB) (BasketRepository, error) {
	basketCreationQuery := `CREATE TABLE IF NOT EXISTS basket_items
	customerLogin TEXT,
	itemID SERIAL,
	name TEXT,
	itemPrice INTEGER,
	amount INTEGER,
	PRIMARY KEY (customerLogin, itemID)`

	_, creationErr := db.Exec(basketCreationQuery)
	if creationErr != nil {
		return BasketRepository{}, fmt.Errorf("Error creating table 'basket_items': \n%w", creationErr)
	}
	return BasketRepository{db}, nil
}

func (br BasketRepository) GetBasket(customerLogin string) ([]model.TaskCreationDTO, error) {
	query := `SELECT * 
	FROM basket_items
	WHERE customerLogin = $1`

	rows, queryErr := br.Conn.Query(query, customerLogin)
	if queryErr != nil {
		return nil, fmt.Errorf("Error executing query \"%s\" to table \"basket_items\":\n %w", query, queryErr)
	}
	tasks := make([]model.TaskCreationDTO, 0)
	for rows.Next() {
		var taskDTO model.TaskCreationDTO
		scanErr := rows.Scan(&taskDTO.ItemID, &taskDTO.Name, &taskDTO.ItemPrice, &taskDTO.Amount)
		if scanErr != nil {
			return nil, fmt.Errorf("Error scanning values from rows: \n%w", scanErr)
		}
		tasks = append(tasks, taskDTO)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil, fmt.Errorf("Error processing rows: \n%w", iterErr)
	}

	return tasks, nil
}

func (br BasketRepository) SaveToBasket(customerLogin string, toSave model.TaskCreationDTO) error {
	selectQuery := `SELECT *
	FROM basket_items
	WHERE customerLogin = $1 AND itemID = $2`
	rows, selectQueryErr := br.Conn.Query(selectQuery, customerLogin, toSave.ItemID)
	if selectQueryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"basket_items\":\n %w", selectQuery, selectQueryErr)
	}
	isPresent := rows.Next()
	saveQuery := ""
	if isPresent {
		saveQuery = `UPDATE basket_items
		SET name = $3,
		itemPrice = $4,
		amount = $5 
		WHERE customerLogin = $1 AND itemID = $2`
	} else {
		saveQuery = `INSERT INTO basket_items 
		VALUES ($1, $2, $3, $4, $5)`
	}

	_, queryErr := br.Conn.Exec(saveQuery,
		customerLogin,
		toSave.ItemID,
		toSave.Name,
		toSave.ItemPrice,
		toSave.Amount)

	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", saveQuery, queryErr)
	}

	return nil
}

func (br BasketRepository) DeleteFromBasket(customerLogin string, itemID int) error {
	query := `DELETE FROM basket_items
	WHERE customerLogin = $1 AND itemID = $2`
	_, queryErr := br.Conn.Exec(query, customerLogin, itemID)
	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}
	return nil
}

func (br BasketRepository) Begin() (*sql.Tx, error) {
	return br.Conn.Begin()
}

func (br BasketRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return br.Conn.BeginTx(ctx, opts)
}
