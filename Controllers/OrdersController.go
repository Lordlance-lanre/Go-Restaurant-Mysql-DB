package Controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Database"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type orderInput struct {
	UserID    uint   `json:"user_id"`
	FoodID    string `json:"food_id"`
	OrderDate string `json:"order_date"`
}

func parseOrderDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(value))
}

func findOrderReferences(userID uint, foodID string) error {
	var user Models.User
	if err := Database.DB.First(&user, userID).Error; err != nil {
		return err
	}

	var food Models.FoodItems
	if err := Database.DB.Where("food_id = ?", foodID).First(&food).Error; err != nil {
		return err
	}
	return nil
}

func CreateOrder(c fiber.Ctx) error {
	var input orderInput
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse order"})
	}
	if input.UserID == 0 || strings.TrimSpace(input.FoodID) == "" || strings.TrimSpace(input.OrderDate) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id, food_id, and order_date are required"})
	}

	orderDate, err := parseOrderDate(input.OrderDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "order_date must use YYYY-MM-DD format"})
	}
	foodID := strings.TrimSpace(input.FoodID)
	if err := findOrderReferences(input.UserID, foodID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User or food not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	order := Models.Order{User_ID: input.UserID, FoodRef: foodID, Order_Date: orderDate}
	if err := Database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	Database.DB.Preload("User").Preload("FoodItems").First(&order, "order_id = ?", order.Order_ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Order created successfully", "order": order})
}

func GetAllOrders(c fiber.Ctx) error {
	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage", "10"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	var totalOrders int64
	if err := Database.DB.Model(&Models.Order{}).Count(&totalOrders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count orders"})
	}

	var orders []Models.Order
	if err := Database.DB.Preload("User").Preload("FoodItems").Limit(recordsPerPage).Offset((page - 1) * recordsPerPage).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Orders retrieved successfully", "orders": orders, "total": totalOrders, "page": page})
}

func GetOrderByID(c fiber.Ctx) error {
	var order Models.Order
	if err := Database.DB.Preload("User").Preload("FoodItems").Where("order_id = ?", c.Params("order_id")).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order retrieved successfully", "order": order})
}

func UpdateOrders(c fiber.Ctx) error {
	var order Models.Order
	orderID := c.Params("order_id")
	if err := Database.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var input struct {
		UserID    *uint   `json:"user_id"`
		FoodID    *string `json:"food_id"`
		OrderDate *string `json:"order_date"`
	}
	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to parse order"})
	}
	if input.UserID == nil && input.FoodID == nil && input.OrderDate == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Provide user_id, food_id, or order_date to update"})
	}

	updates := make(map[string]interface{})
	if input.UserID != nil {
		if *input.UserID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id must be greater than zero"})
		}
		var user Models.User
		if err := Database.DB.First(&user, *input.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		updates["user_id"] = *input.UserID
	}
	if input.FoodID != nil {
		foodID := strings.TrimSpace(*input.FoodID)
		if foodID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "food_id cannot be empty"})
		}
		var food Models.FoodItems
		if err := Database.DB.Where("food_id = ?", foodID).First(&food).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Food not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		updates["food_id"] = foodID
	}
	if input.OrderDate != nil {
		orderDate, err := parseOrderDate(*input.OrderDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "order_date must use YYYY-MM-DD format"})
		}
		updates["order_date"] = orderDate
	}

	if err := Database.DB.Model(&order).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := Database.DB.Preload("User").Preload("FoodItems").Where("order_id = ?", orderID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Order updated but could not be retrieved"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order updated successfully", "order": order})
}

func DeleteOrders(c fiber.Ctx) error {
	var order Models.Order
	orderID := c.Params("order_id")
	if err := Database.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if err := Database.DB.Delete(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order deleted successfully", "order_id": order.Order_ID})
}