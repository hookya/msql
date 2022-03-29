package msql

import (
	"testing"
)

var dataSourceName = "msql_test:msql_test@/msql_test"

func TestOpen(t *testing.T) {
	_, err := Open(dataSourceName)
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func prepareData() {
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	var sqls = []string{
		`DROP TABLE IF EXISTS users;`,
		`CREATE TABLE users ( 
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT, 
			name VARCHAR(20) NOT NULL DEFAULT '',
			is_vip TINYINT(1) NOT NULL DEFAULT 0,
			PRIMARY KEY (id) 
		);`,
		`INSERT INTO users (name) VALUES ("张三"),("李四"),("王五"),("赵六");`,
	}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}

}
