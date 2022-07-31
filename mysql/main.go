package main

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx" // 驱动
)

type User struct {
	ID           uint
	UUID         string
	NickName     string
	AvatarURL    string `json:"avatar_url"`
	PhoneNum     string
	PassWord     string
	IsAuth       bool // 是否认证
	HeadImageID  uint // 头像
	UpdatedMsgAt *time.Time
	WeChatID     uint
	WebToken     string
	Token        string
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

func main() {
	// 建立连接
	db, err := sqlx.Open(`mysql`, `root:1234567890@tcp(127.0.0.1:3306)/xcar?charset=utf8&parseTime=true`)
	log.Println(db, err)
	// 查询
	// Get查询一个
	// Select一个集合
	// 非查询
	// db.Exec() //执行insert update delete

	mode := User{}
	log.Println("---", mode)
	db.Get(&mode, "select * from tab_user")
	log.Println("---", mode)
}
