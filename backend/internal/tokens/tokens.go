package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

type Token struct {
	PlainText string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int    `json:"user_id"`
	Expiry    int64  `json:"expiry"`
	Scope     string `json:"-"`
}

const (
	ScopeAuth = "authentication"
)

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl).Unix(),
		Scope:  scope,
	}

	emptyBytes := make([]byte, 32) // Adjust size as needed
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	// Generate a random token string (this is a placeholder, implement your own logic)
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)

	// Hash the token (this is a placeholder, implement your own hashing logic)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}
