package scan

import (
	"fmt"
	"reflect"
	"strconv"
)

type baseScanner struct {
	dest reflect.Value
}

func (scanner *baseScanner) Scan(src interface{}) error {
	switch scanner.dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", scanner.dest.Kind())
		}
		s := asString(src)
		i64, err := strconv.ParseInt(s, 10, scanner.dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, scanner.dest.Kind(), err)
		}
		scanner.dest.SetInt(i64)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", scanner.dest.Kind())
		}
		s := asString(src)
		u64, err := strconv.ParseUint(s, 10, scanner.dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, scanner.dest.Kind(), err)
		}
		scanner.dest.SetUint(u64)
		return nil
	case reflect.Float32, reflect.Float64:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", scanner.dest.Kind())
		}
		s := asString(src)
		f64, err := strconv.ParseFloat(s, scanner.dest.Type().Bits())
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, scanner.dest.Kind(), err)
		}
		scanner.dest.SetFloat(f64)
		return nil
	case reflect.String:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", scanner.dest.Kind())
		}
		switch v := src.(type) {
		case string:
			scanner.dest.SetString(v)
			return nil
		case []byte:
			scanner.dest.SetString(string(v))
			return nil
		}
	case reflect.Bool:
		if src == nil {
			return fmt.Errorf("converting NULL to %s is unsupported", scanner.dest.Kind())
		}
		s := asString(src)
		b, err := strconv.ParseBool(s)
		if err != nil {
			err = strconvErr(err)
			return fmt.Errorf("converting driver.Value type %T (%q) to a %s: %v", src, s, scanner.dest.Kind(), err)
		}
		scanner.dest.SetBool(b)
		return nil
	}
	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, scanner.dest.Interface())
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}

	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}
