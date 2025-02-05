package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"userCore/src/handlers"
	"userCore/src/middleware"
)

var (
	auth *handlers.AuthHandler
)

func init() {
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
	r.HandleFunc("/index", handlers.Index).
		Methods(http.MethodGet)
	r.HandleFunc("/logout", auth.Logout).
		Methods(http.MethodGet)
	// 启动服务
	err := http.ListenAndServe("localhost:8086", r)
	if err != nil {
		log.Fatal("启动服务失败，err=", err)
	}
}
