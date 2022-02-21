package msql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DbOrTx interface {
	ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(data interface{}, query string, args ...string)
	QueryCtx(ctx context.Context, data interface{}, query string, args ...string) error
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

func Version() string {
	return version
}
