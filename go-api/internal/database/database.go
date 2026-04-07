package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/workspace/go-api/internal/config"
)

// NewDB opens a MySQL connection using the provided configuration.
// It retries connection attempts every 2 seconds for up to 30 seconds
// to handle Docker container startup ordering.
func NewDB(cfg config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var db *sqlx.DB
	var err error

	for i := 0; i < 15; i++ {
		db, err = sqlx.Connect("mysql", dsn)
		if err == nil {
			log.Println("Connected to MySQL")
			return db, nil
		}
		log.Printf("Waiting for MySQL (%d/15): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to MySQL after 30s: %w", err)
}

// RunMigrations reads and executes all .sql files from the given directory.
func RunMigrations(db *sqlx.DB, migrationsDir string) error {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("reading migration files: %w", err)
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", f, err)
		}
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("executing migration %s: %w", f, err)
			}
		}
		log.Printf("Applied migration: %s", filepath.Base(f))
	}

	return nil
}
