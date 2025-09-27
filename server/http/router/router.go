package router

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/common/middleware"
	_ "github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/docs"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/http/api"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("password", passwordStrengthValidator)
		if err != nil {
			panic(err)
		}
	}

	v1UnAuthed := r.Group("/user-ms/v1")
	{
		// swagger router
		v1UnAuthed.GET("/swagger/*any", gs.WrapHandler(
			swaggerFiles.Handler,
			gs.URL("/user-ms/v1/swagger/doc.json"),
		))
		v1UnAuthed.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		v1UnAuthed.POST("/login", api.UserLogin)
		v1UnAuthed.POST("/users", api.Register)
		v1UnAuthed.PUT("/users/activate", api.Validate)
	}
	v1Authed := r.Group("/user-ms/v1")
	{
		v1Authed.Use(middleware.AuthMiddleware())
		v1Authed.POST(("/logout"), api.UserLogout)
	}
	return r
}

// Custom password validation rules
var passwordStrengthValidator validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// atleast one letter
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	// atleast one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	// min length 8
	isValidLength := len(password) >= 8

	return hasLetter && hasDigit && isValidLength
}
