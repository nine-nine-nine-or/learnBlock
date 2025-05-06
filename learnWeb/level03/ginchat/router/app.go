package router

import (
	"ginchat/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	{
		api.GET("index", service.GetIndex)
		api.GET("/users/list", service.GetUserList)
		api.POST("/users/create", service.CreateUser)
		api.GET("/users/delete/:id", service.DeleteUser)
		api.POST("/users/update", service.UpdateUser)
		api.POST("/users/loginIn", service.LoginIn)
		api.GET("/users/sendMessage", service.SendMessage)

	}

	return r
}
