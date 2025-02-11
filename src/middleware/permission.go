package middleware

import (
	"log"
	"net/http"
	"userCore/src/db"
	"userCore/src/utils"
)

type PermissionMiddleWareBuilder struct {
	controlPaths []string
}

func NewPermissionMiddleWareBuilder() *PermissionMiddleWareBuilder {
	return &PermissionMiddleWareBuilder{}
}

func (l *PermissionMiddleWareBuilder) ControlPathsAdd(path string) *PermissionMiddleWareBuilder {
	l.controlPaths = append(l.controlPaths, path)
	return l
}

func (l *PermissionMiddleWareBuilder) Build(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies, _ := r.Cookie("userId")
		if cookies != nil {
			userId := cookies.Value
			user, _ := db.QueryUserById(userId)
			roleId := user.RoleId

			flag := 0
			for _, path := range l.controlPaths {
				if r.URL.Path == path {
					flag = 1
				}
			}

			if flag == 1 {
				// 含有路径，鉴权
				if roleId == 1 {
					log.Println(r.URL.Path, ":权限不足")
					url := utils.GetURL(r, "failure", "权限不足！")
					http.Redirect(w, r, url, 302)
				} else if roleId == 2 {
					log.Println(r.URL.Path, ":权限足够")
					next.ServeHTTP(w, r)
				}
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}
