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

func transform(data interface{}, rows *sql.Rows) error {
	val := reflect.ValueOf(data)
	// typ := reflect.TypeOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Slice:
	case reflect.Struct:
		convertStruct(val, rows)
	case reflect.Map:
	default:
		panic("un support type")
	}
	return nil
}

// 结构体类型
func convertStruct(value reflect.Value, rows *sql.Rows) {
	// for i := 0; i < value.NumField(); i++ {
	// 	if value.Field(i).Kind()
	// }
	rows.Scan()
}

// 基础类型
func convertBase(typ reflect.Value, rows *sql.Rows) {
	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int16, reflect.Int64:

	}
}
