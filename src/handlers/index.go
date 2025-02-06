package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"userCore/src/db"
)

type Count struct {
	AllUserCount        int
	AdministratorsCount int
	OrdinaryUsersCount  int
}

type Data struct {
	Date7   []string `json:"date7"`
	Count7  []string `json:"count7"`
	Date30  []string `json:"date30"`
	Count30 []string `json:"count30"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	c := Count{}
	c.AllUserCount = db.GetAllUserCount()
	c.AdministratorsCount = db.GetAllAdministratorsCount()
	c.OrdinaryUsersCount = db.GetAllOrdinaryUsersCount()
	tmpl, _ := template.ParseFiles("view/index.html")
	tmpl.Execute(w, c)
}

func IndexData(w http.ResponseWriter, r *http.Request) {
	var d Data
	d.Date7, d.Count7 = db.NewUsers(7)
	d.Date30, d.Count30 = db.NewUsers(30)
	bytes, _ := json.Marshal(&d)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	//fmt.Println(string(bytes))
	fmt.Fprintln(w, string(bytes))
}
