package handlers

import (
	"html/template"
	"net/http"
)

func UserList(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("view/userList.html")
	tmpl.Execute(w, nil)
}
