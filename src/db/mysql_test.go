package db

import (
	"fmt"
	"testing"
	"userCore/src/config"
)

func TestQueryByPage(t *testing.T) {
	cf := config.LoadConfig(`C:\Users\ASUS\Desktop\GOGOGO\userCore\config.ini`)
	InitDB(cf)
	users := QueryByPage(0, 10)
	fmt.Println(users)
}
