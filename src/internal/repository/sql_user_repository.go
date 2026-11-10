package repository

import (
	"Notifications/src/internal/domain"
	"Notifications/src/internal/service"
	"context"
	"database/sql"
	"uuid"
)

type SQLUserRepository struct {
	db *sql.DB
}

type UserDB struct {
	ID           string
	Name         string
	LastName     string
	Email        string
	PhoneNumber  string
	DeviceToken  string
	PasswordHash string
	IsActive     bool
}

func (u *UserDB) toDomain() *domain.User {
	return &domain.User{
		ID:           uuid.MustParse(u.ID),
		Name:         u.Name,
		LastName:     u.LastName,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		DeviceToken:  u.DeviceToken,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
	}
}

func fromDomain(u *domain.User) *UserDB {
	return &UserDB{
		ID:           u.ID.String(),
		Name:         u.Name,
		LastName:     u.LastName,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		DeviceToken:  u.DeviceToken,
		PasswordHash: u.PasswordHash,
		IsActive:     u.IsActive,
	}
}

func NewSQLUserRepository(db *sql.DB) service.UserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, last_name, email, phone_number, device_token, password_hash, is_active FROM user WHERE email = $1`, email)

	var u UserDB
	if err := row.Scan(&u.ID, &u.Name, &u.LastName, &u.Email, &u.PhoneNumber, &u.DeviceToken, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return u.toDomain(), nil
}

func (r *SQLUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, last_name, email, phone_number, device_token, password_hash, is_active FROM user WHERE id = $1`, id.String())

	var u UserDB
	if err := row.Scan(&u.ID, &u.Name, &u.LastName, &u.Email, &u.PhoneNumber, &u.DeviceToken, &u.PasswordHash, &u.IsActive); err != nil {
		return nil, err
	}
	return u.toDomain(), nil
}

func (r *SQLUserRepository) Create(ctx context.Context, user *domain.User) error {
	u := fromDomain(user)
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO user (id, name, last_name, email, phone_number, device_token, password_hash, is_active) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		u.ID, u.Name, u.LastName, u.Email, u.PhoneNumber, u.DeviceToken, u.PasswordHash, u.IsActive)
	return err
}
