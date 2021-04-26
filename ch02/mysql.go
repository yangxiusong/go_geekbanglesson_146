package ch02

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDb() {
	if db != nil {
		return
	}
	dsn := "root:root@tcp(127.0.0.1:3306)/test?parseTime=true"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(30)

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}
