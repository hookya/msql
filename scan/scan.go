package scan

import (
	"bytes"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"unicode"

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
		return errors.New("msql: data must be a pointer")
	}
	if ptr.IsNil() {
		return errors.New("msql: data is a nil pointer")
	}
	columns, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return errors.New("msql: no columns")
	}
	target := ptr.Elem()
	switch target.Kind() {
	case reflect.Slice:
		typ := target.Type().Elem()
		for rows.Next() {
			elem := reflect.New(typ).Elem()
			if err := ScanRow(rows, elem); err != nil {
				return err
			}

			target.Set(reflect.Append(target, elem))
		}
	default:
		if rows.Next() {
			if err := ScanRow(rows, target); err != nil {
				return err
			}
		}

	}
	return nil
}

func ScanRow(rows *sql.Rows, target reflect.Value) error {
	switch target.Kind() {
	case reflect.Struct:
		return scan2struct(rows, target)
	case reflect.Map:
		return errors.New("msql: un support scan to map")
	default:
		if err := rows.Scan(&baseScanner{target}); err != nil {
			return err
		}
	}
	return nil
}

func scan2struct(rows *sql.Rows, target reflect.Value) error {
	columns, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	var scanner []interface{}
	for _, column := range columns {
		filedByName(target.Type(), func(i int, name string) bool {
			if column.Name() == name {
				scanner = append(scanner, &baseScanner{target.Field(i)})
				return true
			}
			return false
		})
	}
	return rows.Scan(scanner...)
}

func filedByName(typ reflect.Type, fn func(i int, name string) bool) {
	for i := 0; i < typ.NumField(); i++ {
		tag := typ.Field(i).Tag.Get(`json`)
		if fn(i, Camel2Case(tag)) {
			break
		}
	}
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

// utils func

// Camel2Case 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := bytes.Buffer{}
	var err error
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				if err = buffer.WriteByte('_'); err != nil {
					panic(err)
				}
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

// Case2Camel 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}
