package routes

import (
	controller "kharchasathi.com/controllers"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	router.GET("/users", controller.GetUsers())
	router.GET("/users/:id", controller.GetUserById())
	router.PUT("/users/:id", controller.UpdateUser())
	router.DELETE("/users/:id", controller.DeleteUser())
	router.POST("/users", controller.SignUp())
	router.POST("/users/login", controller.Login())
	router.POST("/users/google-login", controller.LoginWithGoogle())
	router.POST("/users/phone", controller.SetUserPhoneNumber())
	router.PATCH("/users/password", controller.ChangeUserPassword())
	router.PATCH("/users/currency", controller.ChangeInitialCurrency())

	router.PATCH("/users/income", controller.AddIncome())
}
