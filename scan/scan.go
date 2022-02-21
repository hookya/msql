package scan

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/shopspring/decimal"
)

func Scan(rows *sql.Rows, data interface{}) error {
	if scanner := trySqlScanner(data); scanner != nil {
		if rows.Next() {
			if err := rows.Scan(scanner); err != nil {
				return err
			}
		}
		return rows.Err()
	}
	ptr := reflect.ValueOf(data)
	if ptr.Kind() != reflect.Ptr {
		return errors.New("msql: data must be a pointer.")
	}
	if ptr.IsNil() {
		return errors.New("msql: data is a nil pointer.")
	}
	columns, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return errors.New("msql: no columns.")
	}
	return nil
}

func trySqlScanner(ptr interface{}) sql.Scanner {
	if scanner, ok := ptr.(sql.Scanner); ok {
		switch v := ptr.(type) {
		case *decimal.Decimal:
			return decimalScanner{v}
		}
		return scanner
	}
	return nil
}

type decimalScanner struct {
	d *decimal.Decimal
}

func (ds decimalScanner) Scan(src interface{}) error {
	if src == nil {
		*ds.d = decimal.Zero
		return nil
	}
	return ds.d.Scan(src)
}
