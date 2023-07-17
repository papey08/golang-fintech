package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
	"homework8/internal/model/errs"
	"net/http"
)

// Метод для получения пользователя по его ID
func getUserByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getID(c, "user_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		user, getErr := a.GetUserByID(c, userID)

		switch getErr {
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(getErr))
		case nil:
			c.JSON(http.StatusOK, UserSuccessResponse(&user))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(getErr))
		}
	}
}

// Метод для создания пользователя (user)
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		newUser, createErr := a.CreateUser(c, reqBody.ID, reqBody.Nickname, reqBody.Email)

		switch createErr {
		case errs.UserAlreadyExists:
			c.AbortWithStatusJSON(http.StatusConflict, ErrorResponse(createErr))
		case nil:
			c.JSON(http.StatusOK, UserSuccessResponse(&newUser))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(createErr))
		}
	}
}

// Метод для обновления nickname или email пользователя
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		userID, err := getID(c, "user_id")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse(err))
		}

		user, updateErr := a.UpdateUser(c, userID, reqBody.Nickname, reqBody.Email)

		switch updateErr {
		case errs.UserNotExist:
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse(updateErr))
		case nil:
			c.JSON(http.StatusOK, UserSuccessResponse(&user))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse(updateErr))
		}
	}
}
