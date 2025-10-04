package api

import (
	"net/http"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/http/data"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/service"
	"github.com/gin-gonic/gin"
)

// Get Current UserProfile.
// @Summary Get User Profile
// @Description This endpoint allows current login user fetch his/her profile in JSON format.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} data.UserProfileVO
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self [get]
func GetUserProfile(c *gin.Context) {
	userId, exist := c.Get("userID")
	if !exist || userId.(int) <= 0 {
		c.JSON(http.StatusUnauthorized, data.BaseResponse{Code: http.StatusUnauthorized, ErrMsg: "Unauthorized"})
		return
	}
	userProfile, err := service.GetUserProfileService().GetUserProfile(c, userId.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	if userProfile == nil {
		c.JSON(http.StatusNotFound, data.BaseResponse{Code: http.StatusNotFound, ErrMsg: "User not found"})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: userProfile})
}

// Update UserProfile.
// @Summary Update User Profile
// @Description This endpoint allows current login user update his/her profile in JSON format.
// @Tags User
// @Accept json
// @Produce json
// @Param user body data.UserProfileVO true "User profile to update"
// @Success 200 {object} data.UserProfileVO
// @Failure 500 {object} data.BaseResponse
// @Router /user-ms/v1/customer/users/self [put]
func UpdateUserProfile(c *gin.Context) {
	var userProfile data.UserProfileVO
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "Invalid input"})
		return
	}
	if userProfile.ID == 0 {
		c.JSON(http.StatusBadRequest, data.BaseResponse{Code: http.StatusBadRequest, ErrMsg: "User ID is required"})
		return
	}
	userProfile.Email = "" // Email should not be updated here
	err := service.GetUserProfileService().UpdateUserProfile(c, userProfile.ID, &userProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, data.BaseResponse{Code: http.StatusInternalServerError, ErrMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.BaseResponse{Code: http.StatusOK, Data: userProfile})
}
