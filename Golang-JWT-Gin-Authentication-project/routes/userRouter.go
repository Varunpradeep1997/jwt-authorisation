package routes

import(
	controller "github.com/Varunpradeep1997/golang-jwt-projects/controllers"
	"github.com/Varunpradeep1997/golang-jwt-projects/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}