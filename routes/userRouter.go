package routes //this is reached when signup and login is completed from authrouter.

import (
	controller "github.com/Varunpradeep1997/golang-jwt-project/controllers"
	"github.com/Varunpradeep1997/golang-jwt-project/middleware"
	"github.com/gin-gonic/gin"
)

//signup and login anre not protected routes which are in authROUTES
func UserRoutes(incomingRoutes *gin.Engine){               //because untill that token is not provided.
	incomingRoutes.Use(middleware.Authenticate())    
	incomingRoutes.GET("/users",controller.GetUsers())
	incomingRoutes.GET("/users/:user_id",controller.GetUser())
	

}