package msql

import (
	"fmt"
	"testing"
)

func TestTx_Query(t *testing.T) {
	prepareData()
	db, err := Open(dataSourceName)
	if err != nil {
		panic(err)
	}
	var users []User
	if err = db.RunInTrans(func(tx *Tx) error {
		//err = tx.Query(&users, "SELECT * FROM users")
		//if err != nil {
		//	return err
		//}
		_, err = tx.Exec(`INSERT INTO users (name,is_vip) VALUES ("测试一",true),("测试二",false)`)
		if err != nil {
			return err
		}
		err = tx.Query(&users, "SELECT * FROM users WHERE is_vip = true")
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Printf("%v", users)
}
