package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 全局变量
var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// 连接数据库
func initMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func main() {

	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	log.Println("连接数据库成功")

	r := gin.Default()

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
