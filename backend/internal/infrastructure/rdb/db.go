package rdb

import (
	"fmt"
	"strconv"

	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stephenafamo/bob"
)

// NewDB creates a new database connection and returns it as a bob.Executor.
func NewDB(cfg *env.Config) (bob.Executor, func(), error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, strconv.Itoa(cfg.DB.Port), cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	db, err := bob.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping the underlying *sql.DB
	if err := db.DB.Ping(); err != nil {
		db.DB.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	cleanup := func() {
		db.DB.Close()
	}

	return db, cleanup, nil
}
