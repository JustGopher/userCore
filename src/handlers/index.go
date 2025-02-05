package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("访问index")
	tmpl, _ := template.ParseFiles("view/index.html")
	tmpl.Execute(w, nil)
}
