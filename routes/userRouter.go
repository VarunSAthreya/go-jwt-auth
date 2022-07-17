package routes

import (
	"github.com/VarunSAthreya/go-jwt-auth/controllers"
	middleware "github.com/VarunSAthreya/go-jwt-auth/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	r.Use(middleware.Authenticate())
	r.GET("/users", controllers.GetUsers())
	r.GET("/users/:userId", controllers.GetUser())
}
