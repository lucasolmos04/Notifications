package domain

import "uuid"

type User struct {
	ID           uuid.UUID
	Name         string
	LastName     string
	Email        string
	PhoneNumber  string
	DeviceToken  string
	PasswordHash string
	IsActive     bool
}
