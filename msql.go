package msql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	db *sql.DB
}

func Open(dataSourceName string) (*Db, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Db{
		db: db,
	}, nil
}
