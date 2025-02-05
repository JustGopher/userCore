package middleware

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type LoginMiddleWareBuilder struct {
	paths []string
}

func NewLoginMiddleWareBuilder() *LoginMiddleWareBuilder {
	return &LoginMiddleWareBuilder{}
}

func (l *LoginMiddleWareBuilder) IgnorePaths(path string) *LoginMiddleWareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddleWareBuilder) Build(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range l.paths {
			if r.URL.Path == path {
				next.ServeHTTP(w, r)
				return
			}
		}

		cookies, err := r.Cookie("user")
		// cookies 不存在，则重定向
		if err != nil || cookies.Value == "" {
			log.Println(r.URL.Path, ":cookie不存在")
			w.WriteHeader(301)
			tmpl, _ := template.ParseFiles("view/login.html")
			err := tmpl.Execute(w, nil)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		// 存在
		log.Println(r.URL.Path, ":cookie存在,", cookies.Value)
		next.ServeHTTP(w, r)
	})
}
