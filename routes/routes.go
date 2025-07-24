package routes

import (
	"fmt"
	"golangReact/controllers"
	"golangReact/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	fmt.Println("Setting up router...") 

	//initialize gin
	router := gin.Default()
// set up CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

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
	router.GET("/api/getUser/:id", controllers.FindUserById)
	router.GET("/api/insertMillionUser", controllers.InsertMillionUsers)
	router.POST("/api/import-csv", controllers.ImportUsersFromCSV)
	router.POST("/api/import-declare", controllers.ImportDeclareFromCSV)
	return router
}

