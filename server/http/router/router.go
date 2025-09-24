package router

import (
	"github.com/gin-gonic/gin"

	_ "github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/docs"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	// swagger router
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/user-ms/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		// v1.POST("/users", api.UserLogin)
		// v1.PUT("/users/activate", api.Validate)
		// v1.POST("/login", api.UserLogin)
		// v1.POST(("/logout"), api.UserLogout)
	}
	return r
}
