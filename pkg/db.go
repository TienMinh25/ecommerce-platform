package pkg

import (
	"context"
	"database/sql"
)

type Database interface {
	PrepareStatement(ctx context.Context, query string) (PrepareStatement, error)

	BeginTx(ctx context.Context) (Tx, error)

	Close() error
}

type Tx interface {
	PrepareStatement(ctx context.Context, query string) (PrepareStatement, error)

	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}

type PrepareStatement interface {
	// Exec Executes a SQL query (e.g, INSERT, DELETE, UPDATE) with optional args, returning an errors if the operation fails.
	Exec(ctx context.Context, sql string, args ...any) error

	// QueryRow Executes a SQL query expected to return a single row, returning a Row object for result scanning.
	QueryRow(ctx context.Context, sql string, args ...any) Row

	// Query Executes a SQL query expected to return multiple rows, returning a Rows object or an error
	Query(ctx context.Context, sql string, args ...any) (Rows, error)

	// ExecWithResult Executes a SQL query (e.g, INSERT, DELETE, UPDATE) with optional args, returning an errors or sql result
	ExecWithResult(ctx context.Context, sqlStr string, args ...any) (sql.Result, error)
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Scan(dest ...interface{}) error

	Next() bool

	Close()
}
