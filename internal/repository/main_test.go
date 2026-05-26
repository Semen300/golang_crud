package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testDB *sql.DB

func waitForDB(db *sql.DB) error {
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			return nil
		}
		log.Printf("Try %d: Error connecting to DB, Retrying...\n", i)
		time.Sleep(time.Second)
	}
	return fmt.Errorf("db not ready")
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, containerErr := postgres.RunContainer(ctx,
		postgres.WithDatabase("testDB"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"))
	if containerErr != nil {
		log.Fatal(containerErr)
	}
	defer container.Terminate(ctx)

	connString, formatErr := container.ConnectionString(ctx, "sslmode=disable")
	if formatErr != nil {
		log.Fatal(formatErr)
	}

	db, connErr := sql.Open("pgx", connString)
	if connErr != nil {
		log.Fatal(connErr)
	}
	defer db.Close()

	testDB = db

	if err := waitForDB(db); err != nil {
		log.Fatal(err)
	}

	migrateOrders(testDB)
	migrateItems(testDB)
	migrateTasks(testDB)
	migrateUsers(testDB)
	migrateBasketItems(testDB)
	migrateTokens(testDB)

	code := m.Run()
	os.Exit(code)
}
