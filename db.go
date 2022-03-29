package msql

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hookya/msql/scan"
)

type Db struct {
	db      *sql.DB
	timeout time.Duration
}

func (db *Db) Query(data interface{}, query string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db.QueryCtx(ctx, data, query, args...)
}

func (db *Db) QueryCtx(ctx context.Context, data interface{}, query string, args ...interface{}) error {
	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	typ := reflect.TypeOf(data)
	if typ.Kind() != reflect.Ptr {
		panic(`msql-err: data must be a point`)
	}
	// TODO 扫描查询结果到接口中
	return scan.Scan(rows, data)
}

func (db *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db.ExecCtx(ctx, query, args...)
}

func (db *Db) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

func (db *Db) RunInTrans(fn func(*Tx) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db.RunInTransCtx(ctx, fn)
}

func (db *Db) RunInTransCtx(ctx context.Context, fn func(*Tx) error) error {
	_tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	tx := &Tx{tx: _tx, timeout: defaultTimeout}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.tx.Rollback()
			panic(err)
		}
	}()
	if err = fn(tx); err != nil {
		err = tx.tx.Rollback()
		return err
	}
	if err = tx.tx.Commit(); err != nil {
		_ = tx.tx.Rollback()
		return err
	}
	return nil
}

func (db *Db) SetConnMaxIdleTime(d time.Duration) {
	db.db.SetConnMaxIdleTime(d)
}

func (db *Db) SetConnMaxLifetime(d time.Duration) {
	db.db.SetConnMaxLifetime(d)
}

func (db *Db) SetMaxIdleConns(n int) {
	db.db.SetMaxIdleConns(n)
}

func (db *Db) SetMaxOpenConns(n int) {
	db.db.SetMaxOpenConns(n)
}
