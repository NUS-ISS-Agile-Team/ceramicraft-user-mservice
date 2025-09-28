package api

import (
	"net/http"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/http/data"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/service"
	"github.com/gin-gonic/gin"
)

const tokenExpireDuration = 3600 * 24 * 365 // 1 year

// UserLogin handles user login requests.
//
// @Summary User Login
// @Description Authenticates a user with their email and password and returns a token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body data.UserVO true "User login information"
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200	{object} data.BaseResponse{data=string} "Login successful, returns auth token in cookie"
// @Failure 400 {object} data.BaseResponse{data=string}
// @Failure 500 {object} data.BaseResponse{data=string}
// @Router /user-ms/v1/{client}/login [post]
func UserLogin(c *gin.Context) {
	user := &data.UserVO{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := service.GetLoginService().Login(c, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("auth-token", token, tokenExpireDuration, "/", c.Request.Host, true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// UserLogout handles user logout requests.
//
// @Summary User Logout
// @Description invalidates the user's auth token cookie.
// @Tags Authentication
// @Param client path string true "Client identifier" Enums(customer, merchant)
// @Success 200
// @Router /user-ms/v1/{client}/logout [post]
func UserLogout(c *gin.Context) {
	// Invalidate the auth-token cookie by setting its MaxAge to -1
	c.SetCookie("auth-token", "", -1, "/", c.Request.Host, true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
