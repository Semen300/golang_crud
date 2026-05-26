package repository_test

import (
	"crud-go/internal/model"
	"crud-go/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedTokens = []model.RefreshToken{
	{
		TokenID:   "000000000000AAAA110000000000000",
		Login:     "user1",
		TokenHash: "111111111111111",
		ExpiresAt: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		Revoked:   true,
	},
	{
		TokenID:   "000000000000AAAA110000000000001",
		Login:     "user2",
		TokenHash: "111111111111111",
		ExpiresAt: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
		Revoked:   false,
	},
	{
		TokenID:   "000000000000AAAA110000000000002",
		Login:     "user1",
		TokenHash: "111111111111111",
		ExpiresAt: time.Date(2026, time.June, 2, 0, 0, 0, 0, time.UTC),
		Revoked:   false,
	},
}

func migrateTokens(db *sql.DB) {
	_, createErr := db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
	token_id TEXT PRIMARY KEY,
    login TEXT,
    token_hash TEXT,
    expires_at TIMESTAMP,
    revoked BOOLEAN DEFAULT false
	)`)
	if createErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'tokens': \nError creating table: \n%w", createErr))
	}

	_, insertErr := db.Exec(`INSERT INTO tokens (token_id, login, token_hash, expires_at, revoked)
	values ($1, $2, $3, $4, $5),
	($6, $7, $8, $9, $10),
	($11, $12, $13, $14, $15)`,
		"000000000000AAAA110000000000000", "user1", "111111111111111", time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), true,
		"000000000000AAAA110000000000001", "user2", "111111111111111", time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC), false,
		"000000000000AAAA110000000000002", "user1", "111111111111111", time.Date(2026, time.June, 2, 0, 0, 0, 0, time.UTC), false,
	)
	if insertErr != nil {
		log.Fatal(fmt.Errorf("Error executing migration into table 'tokens': \nError adding values: \n%w", insertErr))
	}
}

func resetTokens(repo *repository.TokenRepository) {
	repo.Conn.Exec(`DROP TABLE IF EXISTS tokens`)
	migrateTokens(repo.Conn)
}

func TestNewTokenRepository(t *testing.T) {
	repo, creationErr := repository.NewTokenRepository(testDB)
	assert.NoError(t, creationErr)
	pingErr := repo.Conn.Ping()
	assert.NoError(t, pingErr)
}

func TestSave(t *testing.T) {
	repo, _ := repository.NewTokenRepository(testDB)
	tokenToSave := model.RefreshToken{
		TokenID:   "000000000000AAAA110000000000003",
		Login:     "user3",
		TokenHash: "111111111111111",
		ExpiresAt: time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC),
	}
	defer resetTokens(&repo)
	saveErr := repo.Save(tokenToSave.TokenID, tokenToSave.Login, tokenToSave.TokenHash, tokenToSave.ExpiresAt)
	assert.NoError(t, saveErr)

	savedToken, _ := repo.GetTokenByLogin(tokenToSave.Login)
	assert.Equal(t, tokenToSave, savedToken)
}

func TestGetTokenByLogin(t *testing.T) {
	testLogin := "user1"
	repo, _ := repository.NewTokenRepository(testDB)
	token, getErr := repo.GetTokenByLogin(testLogin)
	assert.NoError(t, getErr)
	assert.Equal(t, expectedTokens[2], token)
}

func TestRevoke(t *testing.T) {
	testLogin := "user1"
	repo, _ := repository.NewTokenRepository(testDB)
	defer resetTokens(&repo)

	repo.Revoke(testLogin)
	token, _ := repo.GetTokenByLogin(testLogin)
	assert.Equal(t, model.RefreshToken{}, token)
}
