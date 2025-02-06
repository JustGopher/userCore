package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"userCore/src/db"
	"userCore/src/object"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("view/userList.html")
	users := db.QueryByPage(0, 10)
	tmpl.Execute(w, users)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := object.User{}
	user.UserId = r.FormValue("userId")
	user.UserName = r.FormValue("userName")
	user.Email = r.FormValue("email")
	user.Role = r.FormValue("role")
	user.Status, _ = strconv.Atoi(r.FormValue("status"))
	ok := db.UpdateUser(user)
	fmt.Fprintf(w, ok)
}
