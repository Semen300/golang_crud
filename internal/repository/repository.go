package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() *sql.DB {
	connString := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_SCEMA"))
	db, openErr := sql.Open("pgx", connString)
	if openErr != nil {
		log.Panic(openErr)
	}
	return db
}

func Close(db *sql.DB) {
	closeErr := db.Close()
	if closeErr != nil {
		log.Panic(closeErr)
	}
}
