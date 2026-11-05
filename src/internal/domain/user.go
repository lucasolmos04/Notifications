package domain

import "uuid"

type User struct {
	ID           uuid.UUID
	Name         string
	LastName     string
	Email        string
	PasswordHash string
	IsActive     bool
}
