package model

import "time"

type RefreshToken struct {
	TokenID   string    `json:"token_id"`
	Login     string    `json:"login"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}
