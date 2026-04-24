package repository

import (
	"crud-go/internal/model"
	"database/sql"
	"fmt"
	"strconv"
)

type CustomerRepository struct {
	Conn *sql.DB
}

func (CustomerRepository) NewCustomerRepository(db *sql.DB) (CustomerRepository, error) {
	query := `CREATE IF NOT EXISTS TABLE customers(
	id SERIAL PRIMARY KEY,
	login TEXT NOT NULL,
	password TEXT NOT NULL,
	fio TEXT,
	number TEXT UNIQUE NOT NULL,
	email TEXT UNIQUE NOT NULL
	)`
	_, queryErr := db.Exec(query)
	if queryErr != nil {
		return CustomerRepository{}, queryErr
	}
	return CustomerRepository{db}, nil
}

func (c CustomerRepository) GetCustomerById(id int) (model.Customer, error) {
	query := `SELECT * FROM customers
	WHERE id = 
	` + strconv.Itoa(id)
	var DBid int
	var customer model.Customer
	queryErr := c.Conn.QueryRow(query).Scan(&DBid, &customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email)
	if queryErr != nil {
		return model.Customer{}, queryErr
	}
	return customer, nil
}

func (c CustomerRepository) GetCustomerByLogin(login string) (model.Customer, error) {
	query := `SELECT * FROM customers
	WHERE login = 
	` + login
	var id int
	var customer model.Customer
	queryErr := c.Conn.QueryRow(query).Scan(&id, &customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email)
	if queryErr != nil {
		return model.Customer{}, queryErr
	}
	return customer, nil
}

func (c CustomerRepository) AddCustomer(customer model.Customer) error {
	query := fmt.Sprintf(`INSERT INTO customers(login, password, fio, number, email)
	values ("%s", "%s", "%s", "%s", "%s")
	`, customer.Login, customer.Password, customer.Fio, customer.Number, customer.Email)
	_, queryErr := c.Conn.Exec(query)
	return queryErr
}

func (c CustomerRepository) UpdateCustomer(customer model.Customer) error {
	query := fmt.Sprintf(
		`UPDATE customers
		SET password = %s,
			fio = %s,
			number = %s,
			email = %s
		WHERE login = %s`,
		customer.Password,
		customer.Fio,
		customer.Number,
		customer.Number,
		customer.Login)
	_, queryErr := c.Conn.Exec(query)
	return queryErr
}
