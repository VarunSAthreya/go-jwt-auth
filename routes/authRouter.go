package routes

import (
	"github.com/VarunSAthreya/go-jwt-auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	r.POST("users/singup", controllers.Signup())

	r.POST("users/login", controllers.Signin())
}
