package msql

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	prepareData()
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO users (name) VALUES ("测试一"),("测试二")`)
	if err != nil {
		panic(err)
	}
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func TestQuery(t *testing.T) {
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	user := User{}
	err = db.Query(&user, fmt.Sprintf(`SELECT * FROM users WHERE name = "%s"`, `张三`))
	if err != nil {
		panic(err)
	}
	fmt.Println(`aaa`)
}
