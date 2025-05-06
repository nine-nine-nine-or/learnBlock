package main

import (
	"ginchat/models"
	"ginchat/utils"
	"log"
)

func main() {
	//db, err := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//
	//// 迁移 schema
	//
	//// Create
	//user := &models.UserBasic{}
	//user.Name = "帅哥"
	//db.Create(user)
	//
	//// Read
	//fmt.Println(db.First(user, 1)) // 根据整型主键查找
	//// Update - 将 product 的 price 更新为 200
	//db.Model(user).Update("Password", 1234)

	//加载配置
	cfg, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	//初始化数据库
	db, err := utils.InitDB(&cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	err = db.AutoMigrate(&models.UserBasic{})
	if err != nil {
		return
	}

	//user := &models.UserBasic{}
	//user.Name = "美女"
	////测试连接
	//db.Create(user)
	//var users []models.UserBasic
	//db.Find(&users)
}
