package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

var expectedOrders = []model.Order{
	{ID: 1, Name: "order1", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager1", WorkerLogin: "worker1", CustomerLogin: "customer1", Status: 0, PriceTotal: 1000},
	{ID: 2, Name: "order2", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager2", WorkerLogin: "worker2", CustomerLogin: "customer2", Status: 0, PriceTotal: 2000},
	{ID: 3, Name: "order3", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager3", WorkerLogin: "worker3", CustomerLogin: "customer3", Status: 0, PriceTotal: 3000},
}

func migrateOrders(db *sql.DB) {
	_, createErr := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	name TEXT,
	deadline DATE,
	managerLogin TEXT,
	workerLogin TEXT,
	customerLogin TEXT,
	status SERIAL,
	price SERIAL
	)`)
	if createErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'orders': \nError creating table: \n%w", createErr))
	}

	_, insertErr := db.Exec(`INSERT into orders (name, deadline, managerLogin, workerLogin, customerLogin, status, price) VALUES
	($1, $2, $3, $4, $5, $6, $7),
	($8, $9, $10, $11, $12, $13, $14),
	($15, $16, $17, $18, $19, $20, $21)`,
		"order1", time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), "manager1", "worker1", "customer1", 0, 1000,
		"order2", time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), "manager2", "worker2", "customer2", 0, 2000,
		"order3", time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), "manager3", "worker3", "customer3", 0, 3000)
	if insertErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'orders': \nError adding values: \n%w", createErr))
	}
}

func resetOrders(r *repository.OrderRepository) {
	r.Conn.Exec(`DROP TABLE IF EXISTS orders`)
	migrateOrders(r.Conn)
	r.CurrentID = 3
}

func TestNewOrderRepository(t *testing.T) {
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	pingErr := repo.Conn.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
}

func TestGetAllOrders(t *testing.T) {
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)
	orders, getErr := repo.GetAllOrders()
	if getErr != nil {
		t.Fatal(getErr)
	}
	if len(orders) != len(expectedOrders) {
		t.Errorf("Result length missmatch: expected %d, got %d", len(expectedOrders), len(orders))
	}
	for i, o := range orders {
		if !reflect.DeepEqual(o, expectedOrders[i]) {
			t.Errorf("Result missmatch: expected %s, got %s", expectedOrders[i].ToString(), o.ToString())
		}
	}
}

func TestGetOrdersByManager(t *testing.T) {
	var rowNum int = 1
	managerLogin := fmt.Sprintf("manager%d", rowNum)
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)
	orders, getErr := repo.GetOrdersByManager(managerLogin)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if len(orders) != 1 {
		t.Errorf("Result length missmatch: expected %d, got %d", 1, len(orders))
	}

	if !reflect.DeepEqual(orders[0], expectedOrders[rowNum-1]) {
		t.Errorf("Result missmatch: expected %s, got %s", expectedOrders[rowNum-1].ToString(), orders[0].ToString())
	}
}

func TestGetOrdersByWorker(t *testing.T) {
	var rowNum int = 2
	workerLogin := fmt.Sprintf("worker%d", rowNum)
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)
	orders, getErr := repo.GetOrdersByWorker(workerLogin)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if len(orders) != 1 {
		t.Errorf("Result length missmatch: expected %d, got %d", 1, len(orders))
	}

	if !reflect.DeepEqual(orders[0], expectedOrders[rowNum-1]) {
		t.Errorf("Result missmatch: expected %s, got %s", expectedOrders[rowNum-1].ToString(), orders[0].ToString())
	}
}

func TestGetOrdersByCustomer(t *testing.T) {
	var rowNum int = 3
	customerLogin := fmt.Sprintf("customer%d", rowNum)
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)
	orders, getErr := repo.GetOrdersByCustomer(customerLogin)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if len(orders) != 1 {
		t.Errorf("Result length missmatch: expected %d, got %d", 1, len(orders))
	}

	if !reflect.DeepEqual(orders[0], expectedOrders[rowNum-1]) {
		t.Errorf("Result missmatch: expected %s, got %s", expectedOrders[rowNum-1].ToString(), orders[0].ToString())
	}
}

func TestGetOrderByID(t *testing.T) {
	var rowNum int = 1
	id := rowNum
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)
	order, getErr := repo.GetOrderById(id)
	if getErr != nil {
		t.Fatal(getErr)
	}

	if !reflect.DeepEqual(order, expectedOrders[rowNum-1]) {
		t.Errorf("Result missmatch: expected %s, got %s", expectedOrders[rowNum-1].ToString(), order.ToString())
	}
}

func TestSaveOrderSave(t *testing.T) {
	orderToSave := model.Order{ID: 4, Name: "order4", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager4", WorkerLogin: "worker4", CustomerLogin: "customer4", Status: 0, PriceTotal: 4000}
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)

	_, saveErr := repo.SaveOrder(orderToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}

	allOrders, _ := repo.GetAllOrders()
	if len(allOrders) != 4 {
		t.Errorf("Result length missmatch: expected %d, got %d", 4, len(allOrders))
	}

	order, _ := repo.GetOrderById(4)
	if !reflect.DeepEqual(order, orderToSave) {
		t.Errorf("Result missmatch: expected %s, got %s", orderToSave.ToString(), order.ToString())
	}
}

func TestSaveOrderUpdate(t *testing.T) {
	orderToSave := model.Order{ID: 3, Name: "order4", Deadline: time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC), ManagerLogin: "manager4", WorkerLogin: "worker4", CustomerLogin: "customer4", Status: 0, PriceTotal: 4000}
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)

	_, saveErr := repo.SaveOrder(orderToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}

	allOrders, _ := repo.GetAllOrders()
	if len(allOrders) != 3 {
		t.Errorf("Result length missmatch: expected %d, got %d", 3, len(allOrders))
	}

	order, _ := repo.GetOrderById(3)
	if !reflect.DeepEqual(order, orderToSave) {
		t.Errorf("Result missmatch: expected %s, got %s", orderToSave.ToString(), order.ToString())
	}
}

func TestDeleteOrder(t *testing.T) {
	id := 3
	repo, repoErr := repository.NewOrderRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	resetOrders(&repo)

	deleteErr := repo.DeleteOrder(id)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}

	order, _ := repo.GetOrderById(id)
	if !reflect.DeepEqual(order, model.Order{}) {
		t.Error("Delete error: row still exists")
	}
}
