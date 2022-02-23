package routes

import(
	controller "github.com/Varunpradeep1997/golang-jwt-projects/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	http.HandleFunc("/",controller.HandleHome)
	http.HandleFunc("/login",controller.HandleLogin)
	http.HandleFunc("/callback",controller.HandleCallback)
	http.ListenAndServe(":8080", nil)
	
}