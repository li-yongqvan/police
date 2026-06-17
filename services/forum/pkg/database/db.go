package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect establishes a connection to PostgreSQL using pgxpool
func Connect(host, port, dbname, user, password string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// RunMigrations runs pending migrations from the migrations directory
func RunMigrations(dbHost, dbPort, dbName, dbUser, dbPassword string) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=schema_forum_migrations",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	migrationsPath := filepath.Join(".", "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Printf("Migrations directory not found at %s, skipping migrations", migrationsPath)
		return nil
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
