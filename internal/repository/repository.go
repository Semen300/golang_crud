package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_SCEMA"))
	db, openErr := sql.Open("pgx", connString)
	if openErr != nil {
		return nil, openErr
	}
	return db, nil
}

func Close(db *sql.DB) error {
	return db.Close()
}
