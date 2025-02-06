package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"userCore/src/config"
	"userCore/src/db"
	"userCore/src/handlers"
	"userCore/src/middleware"
)

var (
	auth *handlers.AuthHandler
	cf   config.Config
)

func init() {
	cf = config.LoadConfig("./config.ini")
	db.InitDB(cf)
	auth = handlers.NewAuthHandler()
}

func main() {
	r := mux.NewRouter()
	// 注册中间件
	r.Use(middleware.NewLoginMiddleWareBuilder().
		IgnorePaths("/login").
		Build)
	// 注册路由
	r.HandleFunc("/login", auth.Login).
		Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/logout", auth.Logout).
		Methods(http.MethodGet)
	r.HandleFunc("/index", handlers.Index).
		Methods(http.MethodGet)
	r.HandleFunc("/indexData", handlers.IndexData).
		Methods(http.MethodGet)
	r.HandleFunc("/userList", handlers.UserList).
		Methods(http.MethodGet)
	r.HandleFunc("/updateUser", handlers.UpdateUser).
		Methods(http.MethodPost)
	// 启动服务
	db.NewUsers(7)
	err := http.ListenAndServe("localhost:8086", r)
	if err != nil {
		log.Fatal("启动服务失败，err=", err)
	}
}
