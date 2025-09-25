package router

import (
	"github.com/gin-gonic/gin"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/common/middleware"
	_ "github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/docs"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/http/api"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// swagger router
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1UnAuthed := r.Group("/user-ms/v1")
	{
		v1UnAuthed.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1UnAuthed.POST("/login", api.UserLogin)
	}
	v1Authed := r.Group("/user-ms/v1")
	{
		v1Authed.Use(middleware.AuthMiddleware())
		v1Authed.POST(("/logout"), api.UserLogout)
	}
	return r
}
