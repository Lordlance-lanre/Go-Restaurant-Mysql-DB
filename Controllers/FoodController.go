package Controllers

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Database"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func CreateFoods(c fiber.Ctx) error {
	var food Models.FoodItems
	if err := c.Bind().Body(&food); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse food"})
	}

	food.Name = strings.TrimSpace(food.Name)
	food.Food_image = strings.TrimSpace(food.Food_image)
	food.Price = math.Round(food.Price*100) / 100
	food.MenuID = strings.TrimSpace(food.MenuID)
	if food.MenuID == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "menu_id is required",
    })
}
	food.Start_Date = strings.TrimSpace(food.Start_Date)
	food.End_Date = strings.TrimSpace(food.End_Date)
	if food.Name == "" || food.Food_image == "" || food.MenuID == "" || food.Price <= 0 || food.Start_Date == "" || food.End_Date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name, price, food_image, menu_id, start_date, and end_date are required",
		})
	}
	if _, err := time.Parse("2006-01-02", food.Start_Date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "start_date must use YYYY-MM-DD format"})
	}
	if _, err := time.Parse("2006-01-02", food.End_Date); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "end_date must use YYYY-MM-DD format"})
	}
	if food.Start_Date > food.End_Date {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "start_date cannot be after end_date"})
	}
	// "menu_id": "a1ca9085-54c9-452c-94d7-7895297c4a60"

	var menu Models.Menu
	if err := Database.DB.Where("menu_id = ?", food.MenuID).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if err := Database.DB.Create(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Food created successfully",
		"food":    food,
	})
}

func GetAllFoods(c fiber.Ctx) error {
	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage", "10"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * recordsPerPage

	var totalFoods int64
	if err := Database.DB.Model(&Models.FoodItems{}).Count(&totalFoods).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count foods"})
	}

	var foods []Models.FoodItems
	if err := Database.DB.Preload("Menu").Limit(recordsPerPage).Offset(offset).Find(&foods).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Foods retrieved successfully",
		"foods":   foods,
		"total":   totalFoods,
		"page":    page,
	})
}

func GetFoodByID(c fiber.Ctx) error {
	var food Models.FoodItems
	foodID := c.Params("food_id")
	if err := Database.DB.Preload("Menu").Where("food_id = ?", foodID).First(&food).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Food retrieved successfully",
		"food":    food,
	})
}

func UpdateFood(c fiber.Ctx) error {
	foodID := c.Params("food_id")
	var food Models.FoodItems
	if err := Database.DB.Where("food_id = ?", foodID).First(&food).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while fetching food"})
	}

	var data struct {
		Name       *string  `json:"name"`
		Price      *float64 `json:"price"`
		FoodImage  *string  `json:"food_image"`
		// MenuID     *string  `json:"menu_id"`
	}
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse body"})
	}

	updates := make(map[string]interface{})
	if data.Name != nil {
		name := strings.TrimSpace(*data.Name)
		if name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name must be a non-empty string"})
		}
		updates["name"] = name
	}
	if data.Price != nil {
		if *data.Price <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "price must be greater than zero"})
		}
		updates["price"] = *data.Price
	}
	if data.FoodImage != nil {
		foodImage := strings.TrimSpace(*data.FoodImage)
		if foodImage == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "food_image must be a non-empty string"})
		}
		updates["food_image"] = foodImage
	}
	// if data.MenuID != nil {
	// 	menuID := strings.TrimSpace(*data.MenuID)
	// 	if menuID == "" {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "menu_id must be a non-empty string"})
	// 	}
	// 	var menu Models.Menu
	// 	if err := Database.DB.Where("menu_id = ?", menuID).First(&menu).Error; err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Menu not found"})
	// 		}
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	updates["menu_id"] = menuID
	// }

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No valid fields to update"})
	}
	if err := Database.DB.Model(&food).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update food"})
	}
	if err := Database.DB.Preload("Menu").Where("food_id = ?", foodID).First(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Food updated but could not be retrieved"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Food updated successfully",
		"food":    food,
	})
}

func DeleteFood(c fiber.Ctx) error {
	foodID := c.Params("food_id")
	var food Models.FoodItems
	if err := Database.DB.Where("food_id = ?", foodID).First(&food).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while fetching food"})
	}
	if err := Database.DB.Delete(&food).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete food"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Food deleted successfully",
		"food_id": food.Food_ID,
	})
}