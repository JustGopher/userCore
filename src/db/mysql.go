package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:xiaomu.303@tcp(182.42.110.229:3306)/sql_test")
	if err != nil {
		log.Println("连接失败", err)
		return
	} else {
		log.Println("数据库连接成功")
	}

}
