package controllers

import "github.com/gin-gonic/gin"

func SetupRoutes(rounter *gin.Engine) {
	rounter.POST("/signup", userController.SignUp)
	rounter.POST("/signin", userController.SignIn)
}
