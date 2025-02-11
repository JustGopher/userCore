package utils

import "net/http"

func GetURL(r *http.Request, status string, message string) string {
	pageNo := r.FormValue("pageNo")
	sUserName := r.FormValue("sUserName")
	sRoleId := r.FormValue("sRoleId")
	sStatus := r.FormValue("sStatus")
	url := "/userList?" +
		"pageNo=" + pageNo +
		"&status=" + status +
		"&message=" + message +
		"&sUserName=" + sUserName +
		"&sRoleId=" + sRoleId +
		"&sStatus=" + sStatus
	return url
}
