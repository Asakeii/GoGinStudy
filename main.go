package main

import (
	"github.com/gin-contrib/cors"
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

func (Todo) TableName() string {
	// 方法
	return "todos"
}

func initMySQL() (err error) {
	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/bubble?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func createTables() (err error) {
	// 创建表
	err = DB.AutoMigrate(&Todo{})
	return err
}

func insertRecord(title *string) (err error) {
	// 插入新数据
	newRecord := Todo{Title: *title, Status: false}
	err = DB.Create(&newRecord).Error
	return err
}

func getAllRecord() ([]Todo, error) {
	var allRecord []Todo
	err := DB.Find(&allRecord).Error
	return allRecord, err
}

func main() {

	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	log.Println("连接数据库成功")

	// 创建表
	err = createTables()
	if err != nil {
		panic(err)
	}
	log.Println("创建表成功")

	r := gin.Default()
	r.Use(cors.Default()) // 简单允许所有跨域请求

	// 定义路由组
	group := r.Group("/v1")
	{
		group.POST("/todo", func(c *gin.Context) {
			// 新增数据
			var input Todo
			err = c.ShouldBindJSON(&input)
			if err != nil {
				panic(err)
			}

			err = insertRecord(&input.Title)
			if err != nil {
				c.JSON(200, gin.H{"err": err.Error()})
				panic(err)
			} else {
				c.JSON(200, input)
			}
		})

		group.GET("/todo", func(c *gin.Context) {
			// 查询所有数据
			var allRecord []Todo
			allRecord, err = getAllRecord()
			if err != nil {
				c.JSON(200, gin.H{"err": err.Error()})
				panic(err)
			} else {
				c.JSON(200, allRecord)
			}

		})

		group.DELETE("/todo/:id", func(c *gin.Context) {
			// 删除指定数据
			deleteID := c.Param("id")
			DB.Where("ID = ?", deleteID).Delete(&Todo{})
		})

		group.PUT("/todo/:id", func(c *gin.Context) {
			// 删除指定数据
			deleteID := c.Param("id")
			var getTodo Todo
			if err = c.BindJSON(&getTodo); err != nil {
				c.JSON(200, gin.H{"err": err.Error()})
				panic(err)
			}
			DB.Where("id == ?", deleteID).Update("status", getTodo.Status)
		})
	}

	err = r.Run(":9000")
	if err != nil {
		panic(err)
	}
}
