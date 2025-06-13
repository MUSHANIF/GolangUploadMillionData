package routes

import (
	"fmt"
	"golangReact/controllers"
	"golangReact/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	fmt.Println("Setting up router...") 

	//initialize gin
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to Golang  API",
		})
	})

	// route register	
	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)
	router.POST("/api/PostUser", middlewares.AuthMiddleware(), controllers.CreateUser)
	router.GET("/api/getUser", controllers.FindUsers)
	router.GET("/api/insertMillionUser", controllers.InsertMillionUsers)
	router.POST("/api/import-csv", controllers.ImportUsersFromCSV)
	return router
}

