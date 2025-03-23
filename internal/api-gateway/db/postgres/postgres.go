package api_gateway_postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type postgres struct {
	db *pgxpool.Pool
}

// todo: inject tracer for distributed tracing
func NewPostgresSQL(lifecycle fx.Lifecycle, manager *env.EnvManager) (pkg.Database, error) {
	dbClient, err := pgxpool.New(context.Background(), fmt.Sprintf("%s/%s", manager.PostgreSQL.PostgresDSN, common.API_GATEWAY_DB))

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool to database %s, error: %w", common.API_GATEWAY_DB, err)
	}

	if err = dbClient.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping to database %s, error: %w", common.API_GATEWAY_DB, err)
	}

	pg := &postgres{
		db: dbClient,
	}

	// manage lifecycle of application, which is used to disconnect to database when application crash or shutdown
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Printf("Postgres connection is connect successfully to database %s", common.API_GATEWAY_DB)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("Close connectio to database %s.....", common.API_GATEWAY_DB)
			pg.db.Close()
			return nil
		},
	})

	return pg, nil
}

// PrepareStatement implements pkg.Database.
func (p *postgres) PrepareStatement(ctx context.Context, query string) (pkg.PrepareStatement, error) {
	panic("using pgx automatically prepare statement, not using it if using pgx")
}

// Exec implements pkg.Database.
func (p *postgres) Exec(ctx context.Context, sql string, args ...any) error {
	// todo: add open telemetry
	_, err := p.db.Exec(ctx, sql, args)

	if err != nil {
		return err
	}

	return nil
}

// ExecWithResult implements pkg.Database.
func (p *postgres) ExecWithResult(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	// todo: add open telemetry
	commandTag, err := p.db.Exec(ctx, sqlStr, args)

	if err != nil {
		return nil, err
	}

	return &pgxResult{
		commandTag: commandTag,
	}, nil
}

// Query implements pkg.Database.
func (p *postgres) Query(ctx context.Context, sql string, args ...any) (pkg.Rows, error) {
	// todo: add open telemetry
	res, err := p.db.Query(ctx, sql, args)

	if err != nil {
		return nil, err
	}

	return &rows{
		rows: res,
	}, nil
}

// QueryRow implements pkg.Database.
func (p *postgres) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	// todo: add open telemetry
	return p.db.QueryRow(ctx, sql, args)
}

// BeginTx implements pkg.Database.
func (p *postgres) BeginTx(ctx context.Context) (pkg.Tx, error) {
	// TODO: future add open telemetry
	transactionPgx, err := p.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		return nil, err
	}

	return &tx{
		transaction: transactionPgx,
	}, nil
}

type tx struct {
	transaction pgx.Tx
}

// Exec implements pkg.Tx.
func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	// todo: add open telemetry
	_, err := t.transaction.Exec(ctx, sql, args)

	if err != nil {
		return err
	}

	return nil
}

// ExecWithResult implements pkg.Tx.
func (t *tx) ExecWithResult(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	// todo: add open telemetry
	commandTag, err := t.transaction.Exec(ctx, sqlStr, args)

	if err != nil {
		return nil, err
	}

	return &pgxResult{
		commandTag: commandTag,
	}, nil
}

// Query implements pkg.Tx.
func (t *tx) Query(ctx context.Context, sql string, args ...any) (pkg.Rows, error) {
	// todo: add open telemetry
	res, err := t.transaction.Query(ctx, sql, args)

	if err != nil {
		return nil, err
	}

	return &rows{
		rows: res,
	}, nil
}

// QueryRow implements pkg.Tx.
func (t *tx) QueryRow(ctx context.Context, sql string, args ...any) pkg.Row {
	// todo: add open telemetry
	return t.transaction.QueryRow(ctx, sql, args)
}

// Commit implements pkg.Tx.
func (t *tx) Commit(ctx context.Context) error {
	// todo: add open telemetry
	return t.transaction.Commit(ctx)
}

// PrepareStatement implements pkg.Tx.
func (t *tx) PrepareStatement(ctx context.Context, query string) (pkg.PrepareStatement, error) {
	panic("using pgx automatically prepare statement, not using it if using pgx")
}

// Rollback implements pkg.Tx.
func (t *tx) Rollback(ctx context.Context) error {
	// TODO: add open telemetry
	return t.transaction.Rollback(ctx)
}

// pgxResult là wrapper implement sql.Result interface
type pgxResult struct {
	commandTag pgconn.CommandTag
}

// LastInsertId implements sql.Result.
// PostgreSQL không hỗ trợ LastInsertId như MySQL, nên phương thức này luôn
// trả về lỗi "không được hỗ trợ"
func (r *pgxResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId is not supported by PostgreSQL, use RETURNING clause instead")
}

// RowsAffected implements sql.Result.
func (r *pgxResult) RowsAffected() (int64, error) {
	// pgx.CommandTag đã có phương thức RowsAffected() trả về int64
	return r.commandTag.RowsAffected(), nil
}

type rows struct {
	rows pgx.Rows
}

// Close implements pkg.Rows.
func (r *rows) Close() {
	r.rows.Close()
}

// Next implements pkg.Rows.
func (r *rows) Next() bool {
	return r.rows.Next()
}

// Scan implements pkg.Rows.
func (r *rows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}
