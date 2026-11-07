package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" //driver "postgres"
)

func NewPostgresConnection(host, port, user, password, dbname string) (*sql.DB, error) {
	// Create connection chain (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Configurations of pool connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Database response before continue?
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
