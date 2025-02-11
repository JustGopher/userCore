package object

type User struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	Password string `json:"-"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
	Role     string `json:"role"`
	RoleId   int    `json:"-"`
}
