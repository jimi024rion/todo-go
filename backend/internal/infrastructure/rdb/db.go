package rdb

import (
	"fmt"
	"strconv"

	"github.com/XSAM/otelsql"
	"github.com/jimi024rion/todo-go/backend/internal/config/env"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stephenafamo/bob"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewDB creates a new database connection and returns it as a bob.DB.
// The underlying *sql.DB is instrumented with OpenTelemetry via otelsql,
// which automatically creates spans for every query and embeds trace context
// into SQL comments (sqlcommenter format) for DB-side log correlation.
func NewDB(cfg *env.Config) (bob.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host, strconv.Itoa(cfg.DB.Port), cfg.DB.User, cfg.DB.Password, cfg.DB.Name)

	sqlDB, err := otelsql.Open("postgres", dsn,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithSQLCommenter(true),
	)
	if err != nil {
		return bob.DB{}, fmt.Errorf("failed to open database: %w", err)
	}

	db := bob.NewDB(sqlDB)

	if err := db.Ping(); err != nil {
		db.Close()
		return bob.DB{}, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
