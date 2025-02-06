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
	"userCore/src/object"
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

// NewUsers 获取 day 天内日期和新增用户数
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
	rows, err := stmt.Query(day)
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

// GetAllUserCount 获取用户总数
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

// GetAllAdministratorsCount 获取管理员总数
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

// GetAllOrdinaryUsersCount 获取普通用户总数
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

// QueryByPage 分页查询
func QueryByPage(page int, num int) []object.User {
	stmt, err := db.Prepare("select u.user_id,u.user_name,u.email,u.status,r.name as role from user u,role r where u.role_id = r.role_id order by u.user_id limit ?,?;")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	rows, err := stmt.Query(page, num)
	if err != nil {
		fmt.Println("查询失败")
		return nil
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	users := []object.User{}
	for rows.Next() {
		user := object.User{}
		rows.Scan(&(user.UserId), &(user.UserName), &(user.Email), &(user.Status), &(user.Role))
		users = append(users, user)
	}
	return users
}

func QueryUserById(id string) (object.User, error) {
	stmt, err := db.Prepare("SELECT user_id, user_name, email, status, role_id from user where user.user_id=?")
	if err != nil {
		log.Println("预处理失败", err)
		return object.User{}, err
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	user := object.User{}
	err = stmt.QueryRow(id).Scan(&(user.UserId), &(user.UserName), &(user.Email), &(user.Status), &(user.RoleId))
	if err != nil {
		log.Println("查询失败", err)
		return object.User{}, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user object.User) string {
	role := user.Role
	if role == "管理员" {
		user.RoleId = 2
	} else {
		user.RoleId = 1
	}
	stmt, _ := db.Prepare("update user set user_name = ?,email = ?,status= ?,role_id= ? where user_id = ?")
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	r, _ := stmt.Exec(user.UserName, user.Email, user.Status, user.RoleId, user.UserId)
	count, _ := r.RowsAffected()
	if count > 0 {
		log.Println("修改成功")
		return "true"
	} else {
		log.Println("修改失败")
		return "false"
	}
}
