package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mahdihagh80/forms/internal/services"
)

var userController UserController

type userSignInData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserController struct {
	userService    services.UserService
	sessionService services.SessionService
}

func SetupUserController(userService services.UserService, sessionService services.SessionService) {
	userController = UserController{
		userService:    userService,
		sessionService: sessionService,
	}
}

func (u UserController) SignUp(c *gin.Context) {
	userData := services.UserData{}
	if err := c.ShouldBindJSON(&userData); err != nil {
		fmt.Println("err : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding data, please check your request body"})
		return
	}
	ctx := context.TODO()
	userId, err := u.userService.Create(ctx, userData)
	fmt.Println("1111111111111 : ", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	sessionsData, err := u.sessionService.Create(context.WithValue(ctx, "userId", userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    url.QueryEscape(sessionsData.SessionId),
		MaxAge:   sessionsData.MaxAge,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: true,
	})

}

func (u UserController) SignIn(c *gin.Context) {
	user := userSignInData{}
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("err : ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error while binding data, please check your request body"})
		return
	}

	ctx := context.TODO()
	userId, err := u.userService.SignIn(ctx, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	sessionsData, err := u.sessionService.Create(context.WithValue(ctx, "userId", userId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "sessionId",
		Value:    url.QueryEscape(sessionsData.SessionId),
		MaxAge:   sessionsData.MaxAge,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
		HttpOnly: true,
	})
}
