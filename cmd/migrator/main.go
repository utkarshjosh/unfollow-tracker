package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/utkarsh/unfollow-tracker/internal/config"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Parse command line flags
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	switch command {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
	case "status":
		if err := migrationStatus(db); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: migrator <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up      Run all pending migrations")
	fmt.Println("  down    Rollback the last migration")
	fmt.Println("  status  Show migration status")
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	return err
}

func migrateUp(db *sql.DB) error {
	migrations, err := getMigrations("up")
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if applied[m.version] {
			continue
		}

		log.Printf("Applying migration: %s", m.version)

		content, err := os.ReadFile(m.path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", m.version, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", m.version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", m.version); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("✅ Applied: %s", m.version)
	}

	log.Println("All migrations applied!")
	return nil
}

func migrateDown(db *sql.DB) error {
	migrations, err := getMigrations("down")
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Find the last applied migration
	var lastApplied *migration
	for i := len(migrations) - 1; i >= 0; i-- {
		version := strings.Replace(migrations[i].version, ".down", ".up", 1)
		if applied[version] {
			lastApplied = &migrations[i]
			break
		}
	}

	if lastApplied == nil {
		log.Println("No migrations to rollback")
		return nil
	}

	log.Printf("Rolling back: %s", lastApplied.version)

	content, err := os.ReadFile(lastApplied.path)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(content)); err != nil {
		tx.Rollback()
		return err
	}

	upVersion := strings.Replace(lastApplied.version, ".down", ".up", 1)
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", upVersion); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("✅ Rolled back: %s", lastApplied.version)
	return nil
}

func migrationStatus(db *sql.DB) error {
	migrations, err := getMigrations("up")
	if err != nil {
		return err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	fmt.Println("\nMigration Status:")
	fmt.Println("─────────────────────────────────")

	for _, m := range migrations {
		status := "❌ Pending"
		if applied[m.version] {
			status = "✅ Applied"
		}
		fmt.Printf("%s  %s\n", status, m.version)
	}

	return nil
}

type migration struct {
	version string
	path    string
}

func getMigrations(direction string) ([]migration, error) {
	pattern := fmt.Sprintf("migrations/*.%s.sql", direction)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var migrations []migration
	for _, f := range files {
		version := filepath.Base(f)
		migrations = append(migrations, migration{
			version: version,
			path:    f,
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}
