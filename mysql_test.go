package msql

import "testing"

var dataSourceName = "user:password@/dbname"

func TestOpen(t *testing.T) {
	_, err := Open(dataSourceName)
	if err != nil {
		t.Errorf("err: %s", err)
	}
}
