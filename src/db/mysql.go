package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
	"userCore/src/config"
)

var db *sql.DB

func InitDB(cf config.Config) {
	fmt.Println(cf.Mysql.Host)
	mysqlConfig := mysql.Config{
		User:                 cf.Mysql.User,
		Passwd:               cf.Mysql.Password,
		Net:                  "tcp",
		Addr:                 cf.Mysql.Host,
		DBName:               cf.Mysql.DBName,
		Timeout:              5 * time.Second,
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         10 * time.Second,
		AllowNativePasswords: true,
		ParseTime:            true, //在查询结果中自动将 MySQL 中的 DATETIME 或 TIMESTAMP 类型的字段解析为 Go 的 time.Time 类型
	}
	dsn := mysqlConfig.FormatDSN() // 格式化连接字符串
	log.Println("DSN:", dsn)

	// 使用 dsn 连接数据库
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	// 检查数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Successfully connected to the database!")
}

func NewUsers(day int) ([]string, []string) {
	// 预处理
	stmt, err := db.Prepare("SELECT DATE(created_at) AS registration_date, COUNT(*) AS user_count FROM user WHERE created_at >= CURDATE() - INTERVAL ? DAY GROUP BY registration_date ORDER BY registration_date;")
	if err != nil {
		fmt.Println("预处理失败")
		return nil, nil
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	// 查询
	rows, err := stmt.Query(day - 1)
	if err != nil {
		fmt.Println("查询失败")
		return nil, nil
	}
	// 获取结果
	type result struct {
		date  []string
		count []string
	}
	res := result{}
	res.count = make([]string, 0)
	res.date = make([]string, 0)
	for rows.Next() {
		var d time.Time
		var c int
		rows.Scan(&d, &c)
		s := d.Format("2006-01-02")
		res.date = append(res.date, s)
		res.count = append(res.count, strconv.Itoa(c))
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	return res.date, res.count
}

func GetAllUserCount() int {
	// 预处理
	stmt, err := db.Prepare("SELECT count(*) from user")
	if err != nil {
		log.Println("预处理失败", err)
		return -1
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	// 查询
	var c int
	err = stmt.QueryRow().Scan(&c)
	if err != nil {
		log.Println("查询失败", err)
		return -1
	}
	fmt.Println(c)
	return c
}

func GetAllAdministratorsCount() int {
	// 预处理
	stmt, err := db.Prepare("SELECT count(*) from user where role_id=2")
	if err != nil {
		log.Println("预处理失败", err)
		return -1
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	// 查询
	var c int
	err = stmt.QueryRow().Scan(&c)
	if err != nil {
		log.Println("查询失败", err)
		return -1
	}
	fmt.Println(c)
	return c
}
func GetAllOrdinaryUsersCount() int {
	// 预处理
	stmt, err := db.Prepare("SELECT count(*) from user where role_id=1")
	if err != nil {
		log.Println("预处理失败", err)
		return -1
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	// 查询
	var c int
	err = stmt.QueryRow().Scan(&c)
	if err != nil {
		log.Println("查询失败", err)
		return -1
	}
	fmt.Println(c)
	return c
}
