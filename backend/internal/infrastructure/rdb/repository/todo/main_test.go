package todo_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jimi024rion/todo-go/backend/internal/testhelper"
	"github.com/stephenafamo/bob"
)

var (
	testSQLDB *sql.DB
	testBobDB bob.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	sqlDB, bobDB, teardown, err := testhelper.NewContainerDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup test DB: %v", err)
	}
	defer teardown()

	testSQLDB = sqlDB
	testBobDB = bobDB

	os.Exit(m.Run())
}
