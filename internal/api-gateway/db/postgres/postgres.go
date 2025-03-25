package api_gateway_postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/third_party/tracing"

	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/env"
	"github.com/TienMinh25/ecommerce-platform/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type postgres struct {
	db     *pgxpool.Pool
	tracer pkg.Tracer
}

func NewPostgresSQL(lifecycle fx.Lifecycle, manager *env.EnvManager, tracer pkg.Tracer) (pkg.Database, error) {
	dbClient, err := pgxpool.New(context.Background(), fmt.Sprintf("%s/%s", manager.PostgreSQL.PostgresDSN, common.API_GATEWAY_DB))

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool to database %s, error: %w", common.API_GATEWAY_DB, err)
	}

	if err = dbClient.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cannot ping to database %s, error: %w", common.API_GATEWAY_DB, err)
	}

	pg := &postgres{
		db:     dbClient,
		tracer: tracer,
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

// BeginTxFunc implements pkg.Database
func (p *postgres) BeginTxFunc(ctx context.Context, options pgx.TxOptions, f func(tx pkg.Tx) error) error {
	transaction, err := p.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// Ghi log lỗi rollback nếu cần, nhưng không return
			_ = transaction.Rollback(ctx)
			// Re-panic sau khi đã cleanup resources
			panic(p)
		}
	}()

	if err = f(transaction); err != nil {
		// Lưu lỗi gốc
		originalErr := err
		// Thử rollback, nhưng vẫn ưu tiên lỗi gốc
		if rbErr := transaction.Rollback(ctx); rbErr != nil {
			// Có thể log lỗi rollback hoặc kết hợp với lỗi gốc
			return fmt.Errorf("execution error: %v, rollback error: %v", originalErr, rbErr)
		}
		return originalErr
	}

	// Commit transaction
	return transaction.Commit(ctx)
}

// Exec implements pkg.Database.
func (p *postgres) Exec(ctx context.Context, sql string, args ...any) error {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "Exec"))
	defer span.End()

	_, err := p.db.Exec(ctx, sql, args)

	if err != nil {
		return err
	}

	return nil
}

// ExecWithResult implements pkg.Database.
func (p *postgres) ExecWithResult(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "ExecWithResult"))
	defer span.End()

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
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "Query"))
	defer span.End()

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
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "QueryRow"))
	defer span.End()

	return p.db.QueryRow(ctx, sql, args)
}

// BeginTx implements pkg.Database.
func (p *postgres) BeginTx(ctx context.Context, options pgx.TxOptions) (pkg.Tx, error) {
	ctx, span := p.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.DBLayer, "BeginTx"))
	defer span.End()

	transactionPgx, err := p.db.BeginTx(ctx, options)

	if err != nil {
		return nil, err
	}

	return &tx{
		transaction: transactionPgx,
		tracer:      p.tracer,
	}, nil
}

type tx struct {
	transaction pgx.Tx
	tracer      pkg.Tracer
}

// Exec implements pkg.Tx.
func (t *tx) Exec(ctx context.Context, sql string, args ...any) error {
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "Exec"))
	defer span.End()

	_, err := t.transaction.Exec(ctx, sql, args)

	if err != nil {
		return err
	}

	return nil
}

// ExecWithResult implements pkg.Tx.
func (t *tx) ExecWithResult(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "ExecWithResult"))
	defer span.End()

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
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "Query"))
	defer span.End()

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
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "QueryRow"))
	defer span.End()

	return t.transaction.QueryRow(ctx, sql, args)
}

// Commit implements pkg.Tx.
func (t *tx) Commit(ctx context.Context) error {
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "Commit"))
	defer span.End()

	return t.transaction.Commit(ctx)
}

// PrepareStatement implements pkg.Tx.
func (t *tx) PrepareStatement(ctx context.Context, query string) (pkg.PrepareStatement, error) {
	panic("using pgx automatically prepare statement, not using it if using pgx")
}

// Rollback implements pkg.Tx.
func (t *tx) Rollback(ctx context.Context) error {
	ctx, span := t.tracer.StartFromContext(ctx, tracing.GetSpanName(tracing.TransactionLayer, "Rollback"))
	defer span.End()

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
