package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

var expectedTasks = []model.Task{
	{Id: 1, Name: "task1", ContractID: 1, ItemID: 1, Amount: 1, Finished: true, Price: 100},
	{Id: 2, Name: "task2", ContractID: 2, ItemID: 2, Amount: 2, Finished: true, Price: 200},
	{Id: 3, Name: "task3", ContractID: 3, ItemID: 3, Amount: 3, Finished: true, Price: 300},
}

func migrateTasks(db *sql.DB) {
	_, createErr := db.Exec(`CREATE TABLE IF NOT EXISTS tasks(
	id SERIAL PRIMARY KEY,
	name TEXT,
	contractID SERIAL,
	itemID SERIAL,
	amount SERIAL,
	finished BOOL,
	price SERIAL)`)
	if createErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'tasks': \nError creating table: \n%w", createErr))
	}

	_, insertErr := db.Exec(`INSERT INTO tasks
	VALUES ($1, $2, $3, $4, $5, $6, $7),
	($8, $9, $10, $11, $12, $13, $14),
	($15, $16, $17, $18, $19, $20, $21)`,
		1, "task1", 1, 1, 1, true, 100,
		2, "task2", 2, 2, 2, true, 200,
		3, "task3", 3, 3, 3, true, 300)
	if insertErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'tasks': \nError adding values: \n%w", createErr))
	}
}

func resetTasks(db *sql.DB) {
	db.Exec(`DROP TABLE IF EXISTS tasks`)
	migrateTasks(db)
}

func getNumberOfTasks(db *sql.DB) int {
	rows, _ := db.Query("SELECT * FROM tasks")
	var numberOfTasks = 0
	for rows.Next() {
		numberOfTasks++
	}
	return numberOfTasks
}

func TestNewTaskRepository(t *testing.T) {
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	pingErr := repo.Conn.Ping()
	if pingErr != nil {
		t.Fatal(pingErr)
	}
}

func TestGetAllTasks(t *testing.T) {
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	tasks, getErr := repo.GetAllTasks()
	if getErr != nil {
		t.Fatal(getErr)
	}
	if len(tasks) != len(expectedTasks) {
		t.Errorf("Result length missmatch: expected %d, got %d", len(expectedTasks), len(tasks))
	}
	for i, task := range tasks {
		if task != expectedTasks[i] {
			t.Errorf("Result missmatch: expected %s, got %s", expectedTasks[i].ToString(), task.ToString())
		}
	}
}

func TestGetTasksByContract(t *testing.T) {
	var rowNum = 1
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	tasks, getErr := repo.GetTasksByContract(rowNum)
	if getErr != nil {
		t.Fatal(getErr)
	}
	if len(tasks) != 1 {
		t.Errorf("Result length missmatch: expected %d, got %d", 1, len(tasks))
	}
	for i, task := range tasks {
		if task != expectedTasks[i] {
			t.Errorf("Result missmatch: expected %s, got %s", expectedTasks[i].ToString(), task.ToString())
		}
	}
}

func TestGetTaskById(t *testing.T) {
	var rowNum = 2
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	task, getErr := repo.GetTaskById(rowNum)
	if getErr != nil {
		t.Fatal(getErr)
	}
	if task != expectedTasks[rowNum-1] {
		t.Errorf("Result missmatch: expected %s, got %s", expectedTasks[rowNum-1].ToString(), task.ToString())
	}
}

func TestSaveTaskSave(t *testing.T) {
	var taskToSave = model.Task{Id: 4, Name: "task4", ContractID: 4, ItemID: 4, Amount: 4, Finished: true, Price: 400}
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	saveErr := repo.SaveTask(taskToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}
	defer resetTasks(repo.Conn)

	numberOfTasks := getNumberOfTasks(repo.Conn)
	if numberOfTasks != 4 {
		t.Errorf("Result length missmatch: expected %d, got %d", 4, numberOfTasks)
	}

	task, _ := repo.GetTaskById(4)
	if task != taskToSave {
		t.Errorf("Result missmatch: expected %s, got %s", taskToSave.ToString(), task.ToString())
	}
}

func TestSaveTaskUpdate(t *testing.T) {
	var taskToSave = model.Task{Id: 3, Name: "task4", ContractID: 4, ItemID: 4, Amount: 4, Finished: true, Price: 400}
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	saveErr := repo.SaveTask(taskToSave)
	if saveErr != nil {
		t.Fatal(saveErr)
	}
	defer resetTasks(repo.Conn)

	numberOfTasks := getNumberOfTasks(repo.Conn)
	if numberOfTasks != 3 {
		t.Errorf("Result length missmatch: expected %d, got %d", 4, numberOfTasks)
	}

	task, _ := repo.GetTaskById(3)
	if task != taskToSave {
		t.Errorf("Result missmatch: expected %s, got %s", taskToSave.ToString(), task.ToString())
	}
}

func TestDeleteTask(t *testing.T) {
	var rowNum = 2
	repo, repoErr := repository.NewTaskRepository(testDB)
	if repoErr != nil {
		t.Fatal(repoErr)
	}
	deleteErr := repo.DeleteTask(rowNum)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
	defer resetTasks(repo.Conn)

	task, _ := repo.GetTaskById(rowNum)
	if task != (model.Task{}) {
		t.Error("Delete error: row still exists")
	}
}
