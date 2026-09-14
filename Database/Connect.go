package Database

import (
	"fmt"
	"os"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	fmt.Println("database connection started")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	dsn := os.Getenv("APP_CONNECTION_URL_STRING")
	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("Failed to connect to database", err)
		return
	}

	DB = db
	fmt.Println("Connected to database")

	if err := DB.AutoMigrate(&Models.User{}, &Models.Menu{}, &Models.FoodItems{}); err != nil {
		fmt.Println("AutoMigrate failed:", err)
	}

}
