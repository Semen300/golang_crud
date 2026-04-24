package repository

import (
	"crud-go/internal/model"
	"database/sql"
)

type ContractRepository struct {
	Conn *sql.DB
}

// NewContractRepository создаёт новый репозиторий для доступа к функционалу контрактов. Также проводит инициализацию таблиц "contracts" и "tasks"
// Принимает указатель на подключение к базе данных
// Возвращает новый экземпляр репозитория
// TODO: разобраться с зоной ответсвенности репозитория (к каким таблицам он имеет доступ),
// применить паттерн "синглтон" для создания одного экземпляра репозитория
func NewContractRepository(db *sql.DB) (ContractRepository, error) {
	contractsCreationQuery := `CREATE TABLE IF NOT EXISTS contracts (
	id SERIAL PRIMARY KEY,
	name TEXT,
	deadline DATE,
	managerLogin TEXT,
	workerLogin TEXT,
	customerLogin TEXT,
	status SERIAL,
	price SERIAL
	)`

	tasksCreationQuery := `CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	name TEXT,
	contractID SERIAL,
	itemID SERIAL,
	amount SERIAL,
	finished BOOL,
	price SERIAL
	)
	`
	_, migrationErr := db.Exec(contractsCreationQuery)
	if migrationErr != nil {
		return ContractRepository{}, migrationErr
	}
	_, migrationErr = db.Exec(tasksCreationQuery)
	if migrationErr != nil {
		return ContractRepository{}, migrationErr
	}

	return ContractRepository{db}, nil
}

// GetAllContracts служит для получения всех заказов, хранящихся в базе данных
// Возвращает список всех заказов и ошибку, если она есть
func (c ContractRepository) GetAllContracts() ([]model.Contract, error) {
	contractQuery := `SELECT c.id, c.name, c.deadline, m.login, m.password, m.fio, w.login, w.password, w.fio, w.supLogin, cu.login, cu.password, cu.fio, cu.number, cu.email, c.status, c.price
	FROM contracts c
	LEFT JOIN managers m ON c.managerLogin = m.login
	LEFT JOIN workers w ON c.workerLogin = w.login
	LEFT JOIN customers cu ON c.customerLogin = cu.login`

	tasksQuery := `SELECT id, name, contractID, itemID, amount, finished, price
	FROM tasks`

	rows, querryErr := c.Conn.Query(contractQuery)
	if querryErr != nil {
		return nil, querryErr
	}
	defer rows.Close()

	var contracts []model.Contract = make([]model.Contract, 0)
	for rows.Next() {
		var contract model.Contract
		var manager model.Manager
		var worker model.Worker
		var customer model.Customer
		rows.Scan(&contract.ID, &contract.Name, &contract.Deadline,
			&manager.Login, &manager.Password, &manager.Fio,
			&worker.Login, &worker.Password, &worker.Fio,
			&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email,
			&contract.Status, &contract.PriseTotal)
		contract.Manager = &manager
		contract.Worker = &worker
		contract.Customer = &customer
		contracts = append(contracts, contract)
	}
	if iterErr := rows.Err(); iterErr != nil {
		return nil, iterErr
	}

	taskRows, taskQueryErr := c.Conn.Query(tasksQuery)
	if taskQueryErr != nil {
		return nil, taskQueryErr
	}
	defer taskRows.Close()

	var tasksMap map[int][]model.Task = make(map[int][]model.Task)
	for taskRows.Next() {
		var task model.Task
		var conID, itemID int
		taskRows.Scan(&task.Id, &task.Name, &conID, &itemID, &task.Amount, &task.Finished, &task.Price)
		tasksMap[conID] = append(tasksMap[conID], task)
	}

	if iterErr := taskRows.Err(); iterErr != nil {
		return nil, iterErr
	}

	for i, contract := range contracts {
		contracts[i].Tasks = tasksMap[contract.ID]
	}

	return contracts, nil
}

// GetContractsByManager служит для получения всех заказов, назначенных определённому менеджеру.
// Принимает логин менеджера, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c ContractRepository) GetContractsByManager(managerLogin string) ([]model.Contract, error) {
	query := `SELECT c.id, c.name, c.deadline, m.login, m.password, m.fio, w.login, w.password, w.fio, w.supLogin, cu.login, cu.password, cu.fio, cu.number, cu.email, c.status, c.price
	FROM contracts c
	LEFT JOIN managers m ON c.managerLogin = m.login
	LEFT JOIN workers w ON c.workerLogin = w.login
	LEFT JOIN customers cu ON c.customerLogin = cu.login
	WHERE c.managerLogin = $1`
	stmt, stmtErr := c.Conn.Prepare(query)
	if stmtErr != nil {
		return nil, stmtErr
	}
	rows, queryErr := stmt.Query(managerLogin)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var contracts []model.Contract = make([]model.Contract, 0)
	for rows.Next() {
		var contract model.Contract
		var manager model.Manager
		var worker model.Worker
		var customer model.Customer
		rows.Scan(&contract.ID, &contract.Name, &contract.Deadline,
			&manager.Login, &manager.Password, &manager.Fio,
			&worker.Login, &worker.Password, &worker.Fio,
			&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email,
			&contract.Status, &contract.PriseTotal)
		contract.Manager = &manager
		contract.Worker = &worker
		contract.Customer = &customer
		contracts = append(contracts, contract)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil, iterErr
	}

	return contracts, nil
}

// GetContractsByWorker служит для получения всех заказов, назначенных определённому работнику.
// Принимает логин работника, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c ContractRepository) GetContractsByWorker(workerLogin string) ([]model.Contract, error) {
	query := `SELECT c.id, c.name, c.deadline, m.login, m.password, m.fio, w.login, w.password, w.fio, w.supLogin, cu.login, cu.password, cu.fio, cu.number, cu.email, c.status, c.price
	FROM contracts c
	LEFT JOIN managers m ON c.managerLogin = m.login
	LEFT JOIN workers w ON c.workerLogin = w.login
	LEFT JOIN customers cu ON c.customerLogin = cu.login
	WHERE c.workerLogin = $1`
	stmt, stmtErr := c.Conn.Prepare(query)
	if stmtErr != nil {
		return nil, stmtErr
	}
	rows, queryErr := stmt.Query(workerLogin)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var contracts []model.Contract = make([]model.Contract, 0)
	for rows.Next() {
		var contract model.Contract
		var manager model.Manager
		var worker model.Worker
		var customer model.Customer
		rows.Scan(&contract.ID, &contract.Name, &contract.Deadline,
			&manager.Login, &manager.Password, &manager.Fio,
			&worker.Login, &worker.Password, &worker.Fio,
			&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email,
			&contract.Status, &contract.PriseTotal)
		contract.Manager = &manager
		contract.Worker = &worker
		contract.Customer = &customer
		contracts = append(contracts, contract)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil, iterErr
	}

	return contracts, nil
}

// GetContractsByCustomer служит для получения всех заказов, оформленных заказчиком.
// Принимает логин заказчика, для которого нужно получить заказы
// Возвращает список заказов и возможную ошибку
func (c ContractRepository) GetContractsByCustomer(customerLogin string) ([]model.Contract, error) {
	query := `SELECT c.id, c.name, c.deadline, m.login, m.password, m.fio, w.login, w.password, w.fio, w.supLogin, cu.login, cu.password, cu.fio, cu.number, cu.email, c.status, c.price
	FROM contracts c
	LEFT JOIN managers m ON c.managerLogin = m.login
	LEFT JOIN workers w ON c.workerLogin = w.login
	LEFT JOIN customers cu ON c.customerLogin = cu.login
	WHERE c.customerLogin = $1`
	stmt, stmtErr := c.Conn.Prepare(query)
	if stmtErr != nil {
		return nil, stmtErr
	}
	rows, queryErr := stmt.Query(customerLogin)
	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var contracts []model.Contract = make([]model.Contract, 0)
	for rows.Next() {
		var contract model.Contract
		var manager model.Manager
		var worker model.Worker
		var customer model.Customer
		rows.Scan(&contract.ID, &contract.Name, &contract.Deadline,
			&manager.Login, &manager.Password, &manager.Fio,
			&worker.Login, &worker.Password, &worker.Fio,
			&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email,
			&contract.Status, &contract.PriseTotal)
		contract.Manager = &manager
		contract.Worker = &worker
		contract.Customer = &customer
		contracts = append(contracts, contract)
	}

	if iterErr := rows.Err(); iterErr != nil {
		return nil, iterErr
	}

	return contracts, nil
}

// GetContractsByID служит для получения заказа по его ID.
// Принимает ID искомого заказа
// Возвращает искомый заказ и возможную ошибку
func (c ContractRepository) GetContractByID(id int) (model.Contract, error) {
	query := `SELECT c.id, c.name, c.deadline, m.login, m.password, m.fio, w.login, w.password, w.fio, w.supLogin, cu.login, cu.password, cu.fio, cu.number, cu.email, c.status, c.price
	FROM contracts c
	LEFT JOIN managers m ON c.managerLogin = m.login
	LEFT JOIN workers w ON c.workerLogin = w.login
	LEFT JOIN customers cu ON c.customerLogin = cu.login
	WHERE c.id = $1`
	var contract model.Contract
	var manager model.Manager
	var worker model.Worker
	var customer model.Customer

	parsingErr := c.Conn.QueryRow(query, id).Scan(&contract.ID, &contract.Name, &contract.Deadline,
		&manager.Login, &manager.Password, &manager.Fio,
		&worker.Login, &worker.Password, &worker.Fio,
		&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email,
		&contract.Status, &contract.PriseTotal)
	if parsingErr != nil {
		return model.Contract{}, parsingErr
	}
	contract.Manager = &manager
	contract.Worker = &worker
	contract.Customer = &customer
	return contract, nil
}

// SaveContract служит для идемпотентного сохранения заказа
// В случае, если сохраняемый заказ существует в БД, обновляет его
// Принимает заказ, который нужно сохранить
// Возвращает возможную ошибку
func (c ContractRepository) SaveContract(contract model.Contract) error {
	selectQuery := `SELECT * 
	FROM contracts
	WHERE id = $1
	`
	rows, queryErr := c.Conn.Query(selectQuery, contract.ID)
	if queryErr != nil {
		return queryErr
	}
	isPresent := rows.Next()
	var saveQuery string
	if !isPresent {
		saveQuery = `INSERT INTO contracts 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	} else {
		saveQuery = `UPDATE contracts 
		SET name=$2,
		deadline=$3,
		managerLogin=$4,
		workerLogin=$5,
		customerLogin=$6,
		status=$7
		WHERE id=$1`
	}

	_, queryErr = c.Conn.Exec(saveQuery,
		contract.ID,
		contract.Name,
		contract.Deadline,
		contract.Manager.Login,
		contract.Worker.Login,
		contract.Customer.Login,
		contract.Status,
		contract.Status,
	)
	if queryErr != nil {
		return queryErr
	}
	return nil

}
