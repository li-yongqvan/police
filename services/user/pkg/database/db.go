package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool   *pgxpool.Pool
	poolMu sync.RWMutex
)

// GetPool returns the current database pool (for middleware access)
func GetPool() (*pgxpool.Pool, error) {
	poolMu.RLock()
	defer poolMu.RUnlock()
	if pool == nil {
		return nil, fmt.Errorf("database pool not initialized")
	}
	return pool, nil
}

// SetPool sets the database pool (used by main.go on startup)
func SetPool(p *pgxpool.Pool) {
	poolMu.Lock()
	defer poolMu.Unlock()
	pool = p
}

// Connect establishes a connection to PostgreSQL using pgxpool
func Connect(host, port, dbname, user, password string) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, dbname, user, password,
	)

	p, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := p.Ping(context.Background()); err != nil {
		p.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return p, nil
}

// RunMigrations runs pending migrations from the migrations directory
func RunMigrations(dbHost, dbPort, dbName, dbUser, dbPassword string) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=schema_auth_migrations",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Look for migrations directory relative to executable
	migrationsPath := filepath.Join(".", "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		// Try relative to source
		migrationsPath = filepath.Join("migrations")
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			log.Printf("Migrations directory not found, skipping migrations")
			return nil
		}
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
