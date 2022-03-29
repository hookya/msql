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
	_, err = db.Exec(`INSERT INTO users (name,is_vip) VALUES ("测试一",true),("测试二",false)`)
	if err != nil {
		panic(err)
	}
}

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	IsVip bool   `json:"isVip"`
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

func TestDb_Query(t *testing.T) {
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	user := User{}
	err = db.Query(&user, fmt.Sprintf(`SELECT * FROM users WHERE id = %d`, 5))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", user)
}

func TestDb_Query2(t *testing.T) {
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	var users []User
	err = db.Query(&users, "SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", users)
}
