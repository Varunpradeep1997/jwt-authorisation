package main

import (
	"log"
	"os"

	routes "github.com/Varunpradeep1997/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	
port := os.Getenv("PORT")

if port == "" {
	port = "8000"
}
router := gin.New()          //here gin is providing router for us 
router.Use(gin.Logger())    //similiar to the gorilla mux router
  

routes.AuthRoutes(router)    //the routes package package is imported and gin is creating router for us.
routes.UserRoutes(router)

router.GET("/api-1", func(c *gin.Context){ //here in gin we need not give w and r like in handlefunction
//gin will take care of that here.
c.JSON(200,gin.H{"SUCCESS":"ACCESS IS GRANTED FOR API-1"})  //bcoz of gin we need not use w.Setheaders here.
})
router.GET("/api-2", func(c *gin.Context){ //here in gin we need not give w and r like in handlefunction
//gin will take care of that here.
c.JSON(200,gin.H{"SUCCESS":"ACCESS IS GRANTED FOR API-2"})  //bcoz of gin we need not use w.Setheaders here.
})

router.Run(":"+ port)
}