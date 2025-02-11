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

type SearchUserList struct {
	UserName string
	RoleId   string
	Status   string
	Page     int
	Num      int
}

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

// QueryUserList 按条件查询用户列表
func QueryUserList(search SearchUserList) []object.User {
	sql1 := `select u.user_id,u.user_name,u.email,u.status,r.name as role from user u,role r where u.role_id = r.role_id `
	sql2 := `order by u.user_id limit ?,?;`
	if search.UserName != "" {
		sql1 += `and u.user_name like "%` + search.UserName + `%" `
	}
	if search.RoleId != "" {
		sql1 += `and u.role_id = ` + search.RoleId + ` `
	}
	if search.Status != "" {
		sql1 += `and u.status = ` + search.Status + ` `
	}
	sql := sql1 + sql2
	log.Println(sql)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	rows, err := stmt.Query((search.Page-1)*search.Num, search.Num)
	if err != nil {
		fmt.Println("查询失败", err)
		return nil
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	users := []object.User{}
	count := 0
	for rows.Next() {
		user := object.User{}
		rows.Scan(&(user.UserId), &(user.UserName), &(user.Email), &(user.Status), &(user.Role))
		users = append(users, user)
		count++
	}
	return users
}

// QueryUserListCount 按条件查询用户数量
func QueryUserListCount(search SearchUserList) int {
	sql1 := `select count(u.user_id) from user u,role r where u.role_id = r.role_id `
	sql2 := `order by u.user_id;`
	if search.UserName != "" {
		sql1 += `and u.user_name like "%` + search.UserName + `%" `
	}
	if search.RoleId != "" {
		sql1 += `and u.role_id = ` + search.RoleId + ` `
	}
	if search.Status != "" {
		sql1 += `and u.status = ` + search.Status + ` `
	}
	sql := sql1 + sql2
	log.Println(sql)
	stmt, err := db.Prepare(sql)
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
	return c
}

// QueryUserById 通过 ID 查询用户
func QueryUserById(id string) (object.User, error) {
	stmt, err := db.Prepare("SELECT user_id, user_name, email, status, role_id from user where user_id=?")
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

// QueryUserByName 通过用户名查询用户
func QueryUserByName(userName string) (object.User, error) {
	stmt, err := db.Prepare("SELECT user_id, password, status from user where user_name=?")
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

	rows, err := stmt.Query(userName)
	if err != nil {
		log.Println("查询失败", err)
		return object.User{}, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if rows.Next() {
		err = rows.Scan(&(user.UserId), &(user.Password), &(user.Status))
		if err != nil {
			return object.User{}, err
		}
	}

	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user object.User) bool {
	stmt, _ := db.Prepare("update user set user_name = ?,email = ?,status= ?,role_id= ? where user_id = ?")
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	r, err := stmt.Exec(user.UserName, user.Email, user.Status, user.RoleId, user.UserId)
	if err != nil {
		log.Println("添加失败", err)
		return false
	}
	count, _ := r.RowsAffected()
	if count > 0 {
		log.Println("修改成功")
		return true
	} else {
		log.Println("修改失败", err)
		return false
	}
}

// UserAdd 添加用户
func UserAdd(user object.User) bool {
	stmt, _ := db.Prepare("insert into user(user_name, password, email, status, role_id) values (?,?,?,?,?)")
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()
	r, err := stmt.Exec(user.UserName, user.Password, user.Email, user.Status, user.RoleId)
	if err != nil {
		log.Println("添加失败", err)
		return false
	}
	count, _ := r.RowsAffected()
	if count > 0 {
		log.Println("添加成功")
		return true
	} else {
		log.Println("添加失败", err)
		return false
	}
}

// UserDelById 删除用户
func UserDelById(roleId int) bool {
	stmt, err := db.Prepare("delete from user where user_id = ?")
	if err != nil {
		log.Println("删除用户失败")
		return false
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	r, err := stmt.Exec(roleId)
	if err != nil {
		log.Println("执行用户失败")
		return false
	}

	count, err := r.RowsAffected()
	if count > 0 {
		log.Println("删除用户成功")
		return true
	} else {
		log.Println("删除用户失败")
		return false
	}
}
