package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stephenafamo/bob"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NewContainerDB はテスト用 PostgreSQL コンテナを起動し、マイグレーション済みの DB を返します。
// TestMain で呼び出してパッケージ内で1回だけ起動してください。
func NewContainerDB(ctx context.Context) (*sql.DB, bob.DB, func(), error) {
	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("todo-go"),
		postgres.WithUsername("user"),
		postgres.WithPassword("P@ssw0rd"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		return nil, bob.DB{}, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, bob.DB{}, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, bob.DB{}, nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := applyMigrations(ctx, sqlDB); err != nil {
		_ = sqlDB.Close()
		_ = container.Terminate(ctx)
		return nil, bob.DB{}, nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	db := bob.NewDB(sqlDB)
	teardown := func() {
		_ = sqlDB.Close()
		_ = container.Terminate(ctx)
	}
	return sqlDB, db, teardown, nil
}

// applyMigrations はマイグレーションファイルを順番に適用します。
func applyMigrations(ctx context.Context, db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	// testhelper/db.go から backend/internal/infrastructure/rdb/migrations/ へのパス
	migrationsDir := filepath.Join(filepath.Dir(filename), "../infrastructure/rdb/migrations")

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, filepath.Join(migrationsDir, e.Name()))
		}
	}
	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", f, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to exec %s: %w", f, err)
		}
	}
	return nil
}

// TruncateTables はすべてのテーブルを空にします。各テストの先頭で呼び出してください。
func TruncateTables(t *testing.T, ctx context.Context, sqlDB *sql.DB) {
	t.Helper()
	_, err := sqlDB.ExecContext(ctx, "TRUNCATE TABLE users CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}
