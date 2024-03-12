package db

import (
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	cfg "github.com/mohitpm/usersvc/config"
)

const (
	migrationFilePath string = "file://db/migrations"
)

func RunMigrations(dbConfig *cfg.Database) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode)

	m, err := migrate.New(migrationFilePath, connStr)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
