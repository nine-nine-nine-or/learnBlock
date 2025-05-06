package main

import (
	_ "ginchat/docs" // 导入生成的docs
	"ginchat/router"
)

// @title           Gin Swagger Demo API
// @version         1.0
// @description     This is a sample server.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	r := router.Router()
	// 添加Swagger路由
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
