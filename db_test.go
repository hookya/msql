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
	//user := User{}
	var id int64 = 0
	err = db.Query(&id, fmt.Sprintf(`SELECT id FROM users WHERE name = "%s"`, `李四`))
	if err != nil {
		panic(err)
	}
	var name string
	err = db.Query(&name, fmt.Sprintf(`SELECT name FROM users WHERE id = %d`, 4))
	if err != nil {
		panic(err)
	}
	var names []string
	err = db.Query(&names, fmt.Sprintf(`SELECT name FROM users`))
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
	fmt.Println(name)
	fmt.Println(names)
}
