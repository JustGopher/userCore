package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"userCore/src/db"
	"userCore/src/object"
	"userCore/src/utils"
)

type UserListData struct {
	Users     []object.User
	Status    string
	Message   string
	UserCount int
	SUserName string
	SRoleId   string
	SStatus   string
	Page      struct {
		PageNo int
		IsHome int
		IsEnd  int
	}
}

func (d *UserListData) SetPage(pageNo int, isHome int, isEnd int) {
	d.Page.PageNo = pageNo
	d.Page.IsHome = isHome
	d.Page.IsEnd = isEnd
}

func UserList(w http.ResponseWriter, r *http.Request) {
	lastPageNo := r.FormValue("pageNo")
	move := r.FormValue("move")
	no, _ := strconv.Atoi(lastPageNo)

	var search db.SearchUserList
	search.UserName = r.FormValue("sUserName")
	search.RoleId = r.FormValue("sRoleId")
	search.Status = r.FormValue("sStatus")
	search.Num = 10

	var userListData = UserListData{}
	userListData.UserCount = db.QueryUserListCount(search)

	endNo := (userListData.UserCount + 10 - 1) / 10
	if lastPageNo == "" {
		search.Page = 1
		if endNo == 1 {
			userListData.SetPage(1, 1, 1)
		} else {
			userListData.SetPage(1, 1, 0)
		}
	} else if lastPageNo != "" && move == "" {
		if no == 1 {
			search.Page = 1
			userListData.SetPage(1, 1, 0)
		} else if no == endNo {
			search.Page = endNo
			userListData.SetPage(endNo, 0, 1)
		} else {
			search.Page = no
			userListData.SetPage(no, 0, 0)
		}
	} else {
		switch move {
		case "home":
			search.Page = 1
			userListData.SetPage(1, 1, 0)
		case "end":
			search.Page = endNo
			userListData.SetPage(endNo, 0, 1)
		case "up":
			if no > 2 {
				search.Page = no - 1
				userListData.SetPage(no-1, 0, 0)
			} else {
				search.Page = no - 1
				userListData.SetPage(no-1, 1, 0)
			}
		case "down":
			if no < endNo-1 {
				search.Page = no + 1
				userListData.SetPage(no+1, 0, 0)
			} else {
				search.Page = endNo
				userListData.SetPage(endNo, 0, 1)
			}
		}
	}
	userListData.Users = db.QueryUserList(search)
	userListData.Message = r.FormValue("message")
	userListData.Status = r.FormValue("status")
	userListData.SUserName = search.UserName
	userListData.SStatus = search.Status
	userListData.SRoleId = search.RoleId

	tmpl, _ := template.ParseFiles("view/userList.html")
	tmpl.Execute(w, userListData)

	userListData.Message = ""
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var message, status string
	user := object.User{}
	user.UserId = r.FormValue("userId")
	user.UserName = r.FormValue("userName")
	user.Email = r.FormValue("email")
	user.RoleId, _ = strconv.Atoi(r.FormValue("roleId"))
	user.Status, _ = strconv.Atoi(r.FormValue("status"))

	ok := db.UpdateUser(user)
	if ok == true {
		message = "修改用户成功！"
		status = "success"
	} else {
		message = "修改用户失败！"
		status = "failure"
	}

	url := utils.GetURL(r, status, message)
	http.Redirect(w, r, url, 302)
}

func UserAdd(w http.ResponseWriter, r *http.Request) {
	var message, status string
	user := object.User{}
	user.UserName = r.FormValue("userName")
	user.Password = r.FormValue("password")
	user.Email = r.FormValue("email")
	user.RoleId, _ = strconv.Atoi(r.FormValue("roleId"))
	user.Status, _ = strconv.Atoi(r.FormValue("status"))

	ok := db.UserAdd(user)
	if ok == true {
		message = "添加用户成功！"
		status = "success"
	} else {
		message = "添加用户失败！"
		status = "failure"
	}

	url := utils.GetURL(r, status, message)
	http.Redirect(w, r, url, 302)
}

func UserDel(w http.ResponseWriter, r *http.Request) {
	var message, status string
	roleId, _ := strconv.Atoi(r.FormValue("userId"))
	fmt.Println(roleId)
	ok := db.UserDelById(roleId)
	if ok == true {
		message = "删除用户成功！"
		status = "success"
	} else {
		message = "删除用户失败！"
		status = "failure"
	}
	url := utils.GetURL(r, status, message)
	http.Redirect(w, r, url, 302)
}
