package infrastructure

import (
	"database/sql"
	"fmt"
	"time"
)

func NewPostgresConnection(host, port, user, password, dbname string) (*sql.DB, error) {
	// Formamos la cadena de conexión (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Configuraciones de optimización del pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Aseguramos que la base de datos realmente responda antes de continuar
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
