package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"userCore/src/db"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type LoginMsg struct {
	Status  string
	Message string
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// GET 登录页面
	if r.Method == http.MethodGet {
		msg := r.FormValue("message")
		if msg == "" {
			//tmpl := template.Must(template.New("login").ParseFiles("view/login.html"))
			tmpl, _ := template.ParseFiles("view/login.html")
			err := tmpl.Execute(w, nil)
			if err != nil {
				log.Fatal("加载页面失败！")
			}
		} else {
			loginMsg := LoginMsg{}
			loginMsg.Status = "failure"
			loginMsg.Message = msg
			//tmpl := template.Must(template.New("login").ParseFiles("view/login.html"))
			tmpl, _ := template.ParseFiles("view/login.html")
			err := tmpl.Execute(w, loginMsg)
			if err != nil {
				log.Fatal("加载页面失败！")
			}
		}
	}

	// POST 登录请求
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println(username, password)
		user, err := db.QueryUserByName(username)
		if err != nil {
			log.Println("查询失败，", err)
			http.Redirect(w, r, "/login?message=系统错误，请重试！", 302)
			return
		}

		// 用户名不存在
		if user.UserId == "" {
			http.Redirect(w, r, "/login?message=用户名或密码错误！", 302)
			return
		}

		// 判断密码是否正确
		if password == user.Password {
			// 判断用户是否禁止登录
			if user.Status == 0 {
				http.Redirect(w, r, "/login?message=该用户禁止登录！", 302)
				return
			} else {
				// 登录成功，设置 Cookie
				http.SetCookie(w, &http.Cookie{
					Name:     "userId",
					Value:    user.UserId,
					Expires:  time.Now().Add(24 * time.Hour),
					HttpOnly: true,
					Path:     "/",
				})
				http.Redirect(w, r, "/index", 302)
				return
			}
		} else {
			http.Redirect(w, r, "/login?message=用户名或密码错误！", 302)
			return
		}
	}
}

// Logout 退出登录
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// 删除 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "user",
		Value:   "",                         // 清空 Cookie 的值
		Expires: time.Now().Add(-time.Hour), // 设置过期时间为过去的时间
		MaxAge:  -1,                         // 使 Cookie 立即过期
	})

	// 重定向到登录页面
	http.Redirect(w, r, "/login", 302)
}
