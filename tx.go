package msql

import (
	"context"
	"database/sql"
	"github.com/hookya/msql/scan"
	"reflect"
	"time"
)

type Tx struct {
	tx      *sql.Tx
	timeout time.Duration
}

func (t *Tx) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	return t.ExecCtx(ctx, query, args...)
}

func (t *Tx) Query(data interface{}, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	return t.QueryCtx(ctx, data, query, args...)
}

func (t *Tx) QueryCtx(ctx context.Context, data interface{}, query string, args ...interface{}) error {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	typ := reflect.TypeOf(data)
	if typ.Kind() != reflect.Ptr {
		panic(`msql-err: data must be a point`)
	}
	return scan.Scan(rows, data)
}
