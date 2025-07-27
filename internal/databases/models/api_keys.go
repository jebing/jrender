package models

import "time"

type ApiKeys struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	APIKey    string     `db:"api_key"`
	ExpiresAt *time.Time `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	Status    int8       `db:"status"`
}
