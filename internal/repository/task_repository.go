package repository

import (
	"crud-go/internal/model"
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	Conn       *sql.DB
	CurrerntID int
}

func NewTaskRepository(db *sql.DB) (TaskRepository, error) {
	query := `CREATE TABLE IF NOT EXISTS tasks(
	id SERIAL PRIMARY KEY,
	name TEXT,
	contractID SERIAL,
	itemID SERIAL,
	amount SERIAL,
	finished BOOL,
	price SERIAL)`

	_, migrationErr := db.Exec(query)
	if migrationErr != nil {
		return TaskRepository{},
			fmt.Errorf("Error creating table \"tasks\":\n %w", migrationErr)
	}
	return TaskRepository{db, 0}, nil
}

func (t TaskRepository) GetTasksByContract(contractID int) ([]model.Task, error) {
	query := `SELECT *
	FROM tasks
	WHERE contractID = $1`

	rows, queryErr := t.Conn.Query(query, contractID)
	if queryErr != nil {
		return nil,
			fmt.Errorf("Error executing query \"%s\" to table \"tasks\":\n %w", query, queryErr)
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		scanErr := rows.Scan(&task.Id, &task.Name, &task.ContractID, &task.ItemID, &task.Amount, &task.Finished, &task.Price)
		if scanErr != nil {
			return nil,
				fmt.Errorf("Error scanning values from rows:\n %w", scanErr)
		}
		tasks = append(tasks, task)
	}
	if iterErr := rows.Err(); iterErr != nil {
		return nil,
			fmt.Errorf("Error processing rows: \n %w", iterErr)
	}

	return tasks, nil
}

func (t TaskRepository) GetTaskById(id int) (model.Task, error) {
	query := `SELECT *
	FROM tasks
	WHERE id = $1`

	var task model.Task
	scanErr := t.Conn.QueryRow(query, id).Scan(&task.Id, &task.Name, &task.ContractID, &task.ItemID, &task.Amount, &task.Finished, &task.Price)
	if scanErr != nil {
		if scanErr == sql.ErrNoRows {
			return model.Task{}, nil
		}
		return model.Task{},
			fmt.Errorf("Error scanning values from rows:\n %w", scanErr)
	}
	return task, nil
}

func (t TaskRepository) SaveTask(task model.Task) error {
	selectQuery := `SELECT *
	FROM tasks
	where id = $1`

	rows, queryErr := t.Conn.Query(selectQuery, task.Id)
	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"tasks\":\n %w", selectQuery, queryErr)
	}
	defer rows.Close()

	isPresent := rows.Next()

	var saveQuery string
	if !isPresent {
		saveQuery = `INSERT INTO tasks
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	} else {
		saveQuery = `UPDATE tasks
		SET name = $2,
		contractID = $3,
		itemID = $4,
		amount = $5,
		finished = $6,
		price = $7
		WHERE id = $1`
	}
	_, queryErr = t.Conn.Exec(saveQuery,
		task.Id,
		task.Name,
		task.ContractID,
		task.ItemID,
		task.Amount,
		task.Finished,
		task.Price)

	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"tasks\":\n %w", saveQuery, queryErr)
	}
	if !isPresent {
		t.CurrerntID++
	}

	return nil
}

func (t TaskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks
	WHERE id = $1`

	_, queryErr := t.Conn.Exec(query, id)
	if queryErr != nil {
		return fmt.Errorf("Error executing query \"%s\" to table \"tasks\":\n %w", query, queryErr)
	}

	return nil
}
