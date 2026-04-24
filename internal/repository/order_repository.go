package repository

import (
	"crud-go/internal/model"
	"database/sql"
	"fmt"
)

// OrderRepository предназначен для выполнения операций, требующих доступа к БД, хранящей список заказов
// TODO: Разобраться с генерацией ID
type OrderRepository struct {
	Conn      *sql.DB // Подключение к БД
	CurrentID uint    // ID последнего заказа в таблице
}

// NewOrderRepository создаёт новый репозиторий для доступа к функционалу контрактов. Также проводит инициализацию таблиц "orders" и "tasks"
// Принимает указатель на подключение к базе данных
// Возвращает новый экземпляр репозитория
// TODO: применить паттерн "синглтон" для создания одного экземпляра репозитория
func NewOrderRepository(db *sql.DB) (OrderRepository, error) {
	ordersCreationQuery := `CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	name TEXT,
	deadline DATE,
	managerLogin TEXT,
	workerLogin TEXT,
	customerLogin TEXT,
	status SERIAL,
	price SERIAL
	)`
	_, migrationErr := db.Exec(ordersCreationQuery)
	if migrationErr != nil {
		return OrderRepository{},
			fmt.Errorf("Error creating table \"orders\":\n %w", migrationErr)
	}
	return OrderRepository{db, 0}, nil
}

// GetAllOrders служит для получения всех заказов, хранящихся в базе данных
// Возвращает список всех заказов и ошибку, если она есть
func (c OrderRepository) GetAllOrders() ([]model.Order, error) {
	orderQuery := `SELECT o.id, o.name, o.deadline, o.managerLogin, o.workerLogin, o.customerLogin, o.status, o.price
	FROM orders o`

	rows, querryErr := c.Conn.Query(orderQuery)
	if querryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", orderQuery, querryErr)
	}
	defer rows.Close()

	var orders []model.Order = make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		scanError := rows.Scan(&order.ID,
			&order.Name,
			&order.Deadline,
			&order.ManagerLogin,
			&order.WorkerLogin,
			&order.CustomerLogin,
			&order.Status,
			&order.PriseTotal)
		if scanError != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanError)
		}
		orders = append(orders, order)
	}
	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}
	return orders, nil
}

// GetOrdersByManager служит для получения всех заказов, назначенных определённому менеджеру.
// Принимает логин менеджера, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c OrderRepository) GetOrdersByManager(managerLogin string) ([]model.Order, error) {
	query := `SELECT o.id, o.name, o.deadline, o.managerLogin, o.workerLogin, o.customerLogin, o.status, o.price
	FROM orders o
	WHERE o.managerLogin = $1`

	rows, queryErr := c.Conn.Query(query, managerLogin)
	if queryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}
	defer rows.Close()

	var orders []model.Order = make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		scanError := rows.Scan(&order.ID,
			&order.Name,
			&order.Deadline,
			&order.ManagerLogin,
			&order.WorkerLogin,
			&order.CustomerLogin,
			&order.Status,
			&order.PriseTotal)
		if scanError != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanError)
		}
		orders = append(orders, order)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}

	return orders, nil
}

// GetOrdersByWorker служит для получения всех заказов, назначенных определённому работнику.
// Принимает логин работника, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c OrderRepository) GetOrdersByWorker(workerLogin string) ([]model.Order, error) {
	query := `SELECT o.id, o.name, o.deadline, o.managerLogin, o.workerLogin, o.customerLogin, o.status, o.price
	FROM orders o
	WHERE o.workerLogin = $1`

	rows, queryErr := c.Conn.Query(query, workerLogin)
	if queryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}
	defer rows.Close()

	var orders []model.Order = make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		scanError := rows.Scan(&order.ID,
			&order.Name,
			&order.Deadline,
			&order.ManagerLogin,
			&order.WorkerLogin,
			&order.CustomerLogin,
			&order.Status,
			&order.PriseTotal)
		if scanError != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanError)
		}
		orders = append(orders, order)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}

	return orders, nil
}

// GetOrdersByCustomer служит для получения всех заказов, оформленных заказчиком.
// Принимает логин заказчика, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c OrderRepository) GetOrdersByCustomer(customerLogin string) ([]model.Order, error) {
	query := `SELECT o.id, o.name, o.deadline, o.managerLogin, o.workerLogin, o.customerLogin, o.status, o.price
	FROM orders o
	WHERE o.customerLogin = $1`

	rows, queryErr := c.Conn.Query(query, customerLogin)
	if queryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}
	defer rows.Close()

	var orders []model.Order = make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		scanError := rows.Scan(&order.ID,
			&order.Name,
			&order.Deadline,
			&order.ManagerLogin,
			&order.WorkerLogin,
			&order.CustomerLogin,
			&order.Status,
			&order.PriseTotal)
		if scanError != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanError)
		}
		orders = append(orders, order)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}

	return orders, nil
}

// GetOrderById служит для получения заказа по его ID.
//
// Принимает ID искомого заказа.
// Возвращает искомый заказ и возможную ошибку.
func (c OrderRepository) GetOrderById(id int) (model.Order, error) {
	query := `SELECT o.id, o.name, o.deadline, o.managerLogin, o.workerLogin, o.customerLogin, o.status, o.price
	FROM orders o
	WHERE o.id = $1`
	var order model.Order

	scanError := c.Conn.QueryRow(query, id).Scan(&order.ID, &order.Name, &order.Deadline,
		&order.ManagerLogin,
		&order.WorkerLogin,
		&order.CustomerLogin,
		&order.Status,
		&order.PriseTotal)
	if scanError != nil {
		return model.Order{},
			fmt.Errorf("Error scanning values from rows:\n %w", scanError)
	}
	return order, nil
}

// SaveOrder служит для идемпотентного сохранения заказа
// В случае, если сохраняемый заказ существует в БД, обновляет его
// Принимает заказ, который нужно сохранить
// Возвращает возможную ошибку
func (o *OrderRepository) SaveOrder(order model.Order) error {
	selectQuery := `SELECT * 
	FROM orders
	WHERE id = $1
	`
	rows, queryErr := o.Conn.Query(selectQuery, order.ID)
	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", selectQuery, queryErr)
	}
	isPresent := rows.Next()
	var saveQuery string
	if !isPresent {
		saveQuery = `INSERT INTO orders 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	} else {
		saveQuery = `UPDATE orders 
		SET name=$2,
		deadline=$3,
		managerLogin=$4,
		workerLogin=$5,
		customerLogin=$6,
		status=$7
		WHERE id=$1`
	}

	_, queryErr = o.Conn.Exec(saveQuery,
		order.ID,
		order.Name,
		order.Deadline,
		order.ManagerLogin,
		order.WorkerLogin,
		order.CustomerLogin,
		order.Status,
		order.Status,
	)
	if queryErr != nil {
		return fmt.Errorf("Error executin Идентификатор g query \"%s\" to table \"orders\":\n %w", saveQuery, queryErr)
	}
	if !isPresent {
		o.CurrentID++
	}
	return nil
}

func (o *OrderRepository) deleteOrder(id string) error {
	query := `DELETE FROM orders
	WHERE id = $1`

	_, queryErr := o.Conn.Exec(query, id)
	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"orders\":\n %w", query, queryErr)
	}

	return nil
}
