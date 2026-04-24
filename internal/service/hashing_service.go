package service

import (
	"crypto/sha256"
	"encoding/hex"
)

func CreateHash(data string) string {
	hashed := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hashed[:])
}
