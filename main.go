package main

import (
	"fmt"
	"golangReact/config"
	"golangReact/database"
	"golangReact/routes"
)

func main() {
	// load config .env
	config.LoadEnv()
	fmt.Println("ENV loaded")

	// inisialisasi database
	database.InitDB()
	fmt.Println("DB connected")

	// gunakan SetupRouter dari routes
	r := routes.SetupRouter()
	fmt.Println("Router setup complete")

	port := config.GetEnv("APP_PORT", "3000")
	fmt.Println("Starting server on port:", port)

	r.Run(":" + port)
}
