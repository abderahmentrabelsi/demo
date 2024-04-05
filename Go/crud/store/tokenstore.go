package store

import (
	"sync"
	"time"
)

var (
	revokedTokens sync.Map
	refreshTokens sync.Map
)

// RevokeToken marks a token as revoked
func RevokeToken(token string) {
	revokedTokens.Store(token, time.Now())
}

// IsTokenRevoked checks if a token is revoked
func IsTokenRevoked(token string) bool {
	_, exists := revokedTokens.Load(token)
	return exists
}

// SetRefreshToken stores a refresh token
func SetRefreshToken(token, email string) {
	refreshTokens.Store(token, email)
}

// GetEmailByRefreshToken retrieves the email associated with a refresh token
func GetEmailByRefreshToken(token string) (string, bool) {
	email, exists := refreshTokens.Load(token)
	if exists {
		return email.(string), true
	}
	return "", false
}

// RemoveRefreshToken deletes a refresh token
func RemoveRefreshToken(token string) {
	refreshTokens.Delete(token)
}
