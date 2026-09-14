package main

import (
	"fmt"
	"os"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Database"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Routes"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	Database.ConnectDB()
	port := os.Getenv("PORT")
	app := fiber.New()
	Routes.AppRoutes(app)
	app.Listen(":" + port)

	fmt.Println("Port:", port)
	fmt.Println("Hello World")
}
