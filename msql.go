package msql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	db      *sql.DB
	timeout time.Duration
}

func Open(dataSourceName string) (*Db, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Db{
		db:      db,
		timeout: defaultTimeout,
	}, nil
}

func (db *Db) Query(data interface{}, query string, args ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db.QueryCtx(ctx, data, query, args...)
}

func (db *Db) QueryCtx(ctx context.Context, data interface{}, query string, args ...string) error {
	rows, err := db.db.QueryContext(ctx, query, args)
	if err != nil {
		return err
	}
	// TODO 扫描查询结果到接口中
	fmt.Println(rows)
	return nil
}

func (db *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timeout)
	defer cancel()
	return db.ExecCtx(ctx, query, args...)
}

func (db *Db) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

func Version() string {
	return version
}
