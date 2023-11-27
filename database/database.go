package database

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Ethea2/nat-server/utils"
)

var DataBase *sql.DB

func ConnectDB() error {
	utils.LoadEnv()
	var err error
	DataBase, err = sql.Open(("mysql"), os.Getenv("DSN"))
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	defer DataBase.Close()
}
