package repository_test

import (
	"crud-go/internal/repository"
	"testing"

	"github.com/joho/godotenv"
)

func TestDBConnection(t *testing.T) {
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		t.Fatal(envErr)
	}
	db, openErr := repository.Connect()
	if openErr != nil {
		t.Fatal(openErr)
	}
	defer db.Close()
	pingErr := db.Ping()
	if pingErr != nil {
		t.Error(pingErr)
	}
}
