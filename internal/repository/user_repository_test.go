package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

var (
	expectedManagers = []model.Manager{
		model.NewManager("manager1", "1111", "Андреев А.А."),
		model.NewManager("manager2", "2222", "Белов Б.Б."),
	}
	expectedWorkers = []model.Worker{
		model.NewWorker("worker1", "1111", "Алексеев А.А.", "manager1"),
		model.NewWorker("worker2", "2222", "Бурый Б.Б.", "mamager2"),
	}
	expectedCustomers = []model.Customer{
		model.NewCustomer("customer1", "1111", "Алёхин А.А.", "89111111111", "alyohin@mail.ru"),
		model.NewCustomer("customer2", "2222", "Банин Б.Б.", "89222222222", "press@rkn.gov.ru"),
	}
)

func migrateUsers(db *sql.DB) {
	_, managersErr := db.Exec(`CREATE TABLE IF NOT EXISTS managers(
	login TEXT PRIMARY KEY,
	password TEXT,
	fio TEXT
	)`)
	if managersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'managers': \nError creating table: \n%w", managersErr))
	}

	_, workersErr := db.Exec(`CREATE TABLE IF NOT EXISTS workers (
	login TEXT PRIMARY KEY,
	password TEXT,
	fio TEXT,
	superiorLogin TEXT
	)`)
	if workersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'workers': \nError creating table: \n%w", workersErr))
	}

	_, customersErr := db.Exec(`CREATE TABLE IF NOT EXISTS customers (
	login TEXT PRIMARY KEY,
	password TEXT,
	fio TEXT,
	number TEXT,
	email TEXT
	)`)
	if customersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'customers': \nError creating table: \n%w", customersErr))
	}

	_, managersErr = db.Exec(`INSERT INTO managers VALUES
	($1, $2, $3),
	($4, $5, $6)`,
		"manager1", "1111", "Андреев А.А.",
		"manager2", "2222", "Белов Б.Б.")

	if managersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'managers': \nError adding values: \n%w", managersErr))
	}

	_, workersErr = db.Exec(`INSERT INTO workers VALUES
	($1, $2, $3, $4),
	($5, $6, $7, $8)`,
		"worker1", "1111", "Алексеев А.А.", "manager1",
		"worker2", "2222", "Бурый Б.Б.", "mamager2")
	if workersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'workers': \nError adding values: \n%w", workersErr))
	}

	_, customersErr = db.Exec(`INSERT INTO customers VALUES
	($1, $2, $3, $4, $5),
	($6, $7, $8, $9, $10)`,
		"customer1", "1111", "Алёхин А.А.", "89111111111", "alyohin@mail.ru",
		"customer2", "2222", "Банин Б.Б.", "89222222222", "press@rkn.gov.ru")
	if customersErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'customers': \nError adding values: \n%w", customersErr))
	}

}

func resetUsers(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS managers")
	db.Exec("DROP TABLE IF EXISTS workers")
	db.Exec("DROP TABLE IF EXISTS customers")
	migrateUsers(db)
}

func getNumberOfCustomers(db *sql.DB) int {
	rows, _ := db.Query("SELECT * FROM customers")
	var numberOfCustomers = 0
	for rows.Next() {
		numberOfCustomers++
	}
	return numberOfCustomers
}

func TestNewUserRepository(t *testing.T) {
	repo, repoErr := repository.NewUserRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	pingErr := repo.Conn.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
}

func TestGetWorkersByManager(t *testing.T) {
	var rowNumber = 1
	repo, repoErr := repository.NewUserRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	managerLogin := expectedManagers[rowNumber-1].Login
	workers, getErr := repo.GetWorkersByManager(managerLogin)
	if getErr != nil {
		t.Fatal(getErr)
	}
	if len(workers) != 1 {
		t.Errorf("Result length missmatch: expected %d, got %d", 1, len(workers))
	}
	if workers[0] != expectedWorkers[rowNumber-1] {
		t.Errorf("Result missmatch: expected %s, got %s", expectedWorkers[rowNumber-1].ToString(), workers[0].ToString())
	}
}

func TestGetRoleByLogin(t *testing.T) {
	login := "customer1"
	expectedPassword := "1111"
	expectedRole := 1
	repo, repoErr := repository.NewUserRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	role, password, getErr := repo.GetRoleByLogin(login)
	if getErr != nil {
		t.Fatal(repoErr)
	}
	if role != expectedRole {
		t.Errorf("Result missmatch: expected %d, got %d", expectedRole, role)
	}
	if password != expectedPassword {
		t.Errorf("Result missmatch: expected %s, got %s", expectedPassword, password)
	}
}

func TestSaveCustomerSave(t *testing.T) {
	customerToSave := model.NewCustomer("customer3", "3333", "Вилкин В.В.", "8933333333", "vilkin@mail.ru")
	repo, repoErr := repository.NewUserRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	saveErr := repo.SaveCustomer(customerToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}
	defer resetUsers(testDB)

	numberOfCustomers := getNumberOfCustomers(repo.Conn)
	if numberOfCustomers != 3 {
		t.Errorf("Result length missmatch: expected %d, got %d", 3, numberOfCustomers)
	}
}

func TestSaveCustomerUpdate(t *testing.T) {
	customerToSave := model.NewCustomer("customer2", "3333", "Вилкин В.В.", "8933333333", "vilkin@mail.ru")
	repo, repoErr := repository.NewUserRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	saveErr := repo.SaveCustomer(customerToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}
	defer resetUsers(testDB)

	numberOfCustomers := getNumberOfCustomers(repo.Conn)
	if numberOfCustomers != 2 {
		t.Errorf("Result length missmatch: expected %d, got %d", 3, numberOfCustomers)
	}

	var customer model.Customer
	queryErr := repo.Conn.QueryRow("SELECT * FROM customers WHERE login = $1", customerToSave.Login).Scan(&customer.Login, &customer.Password, &customer.Fio, &customer.Number, &customer.Email)
	if queryErr != nil {
		t.Errorf("No such row in DB: \n%v", queryErr)
	}

	role, password, passErr := repo.GetRoleByLogin(customerToSave.Login)
	if passErr != nil {
		t.Fatal(passErr)
	}
	log.Println(role, password, passErr)

	if password != customerToSave.Password {
		t.Errorf("Result missmatch: expected %s, got %s", customerToSave.Password, password)
	}

	if customer != customerToSave {
		t.Errorf("Result missmatch: expected %s, got %s", customer.ToString(), customer.ToString())
	}
}
