package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// GET 登录页面
	if r.Method == http.MethodGet {
		//tmpl := template.Must(template.New("login").ParseFiles("login.html"))
		tmpl, _ := template.ParseFiles("view/login.html")
		err := tmpl.Execute(w, nil)
		if err != nil {
			log.Fatal("加载页面失败！")
		}
	}

	// POST 登录请求
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)
		if username == "admin" && password == "123456" {
			// 登录成功，设置 Cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "user",
				Value:    username,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			})
			http.Redirect(w, r, "/index", 302)
			return
			//tmpl, _ := template.ParseFiles("view/index.html")
			//err := tmpl.Execute(w, nil)
			//if err != nil {
			//	log.Fatal("加载页面失败！")
			//}
			//w.WriteHeader(200)
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}
}

// Logout 退出登录
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 删除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user",
		Value:    "",                         // 清空 Cookie 的值
		Expires:  time.Now().Add(-time.Hour), // 设置过期时间为过去的时间
		MaxAge:   -1,                         // 使 Cookie 立即过期
		HttpOnly: true,                       // 使 Cookie 只能通过 HTTP 访问，防止客户端脚本访问
		Path:     "/",                        // 指定路径为根路径，确保删除所有路径下的 session_id
	})

	// 重定向到登录页面
	http.Redirect(w, r, "/login", 302)
}
