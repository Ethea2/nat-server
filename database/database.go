package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DataBase *sql.DB

func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load ENV", err)
		return err
	}

	DataBase, err = sql.Open(("mysql"), os.Getenv("DSN"))
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	defer DataBase.Close()
}
