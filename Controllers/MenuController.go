package Controllers

import (
	// "fmt"
	"errors"
	"strconv"
	"strings"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Database"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func CreateMenu(c fiber.Ctx) error {
	var menu Models.Menu

	if err := c.Bind().Body(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse menu",
		})
	}
	if strings.TrimSpace(menu.Name) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu name is required",
		})
	}
	if strings.TrimSpace(menu.Category) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu category is required",
		})
	}
	if strings.TrimSpace(menu.Start_Date) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu start date is required",
		})
	}
	if strings.TrimSpace(menu.End_Date) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Menu end date is required",
		})
	}

	if err := Database.DB.Create(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu created successfully",
		"menu":    menu,
	})
}

func GetAllMenus(c fiber.Ctx) error {
	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage", "10"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * recordsPerPage

	var totalMenus int64
	if err := Database.DB.Model(&Models.Menu{}).Count(&totalMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count menus",
		})
	}

	var menus []Models.Menu
	if err := Database.DB.Limit(recordsPerPage).Offset(offset).Find(&menus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menus retrieved successfully",
		"menus":   menus,
	})
}

func GetMenuByID(c fiber.Ctx) error {
	var menu Models.Menu
	menuID := c.Params("menu_id")
	if err := Database.DB.Where("menu_id = ?", menuID).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu retrieved successfully",
		"menu":    menu,
	})
}

func UpdateMenu(c fiber.Ctx) error {
	menuID := c.Params("menu_id")
	var menu Models.Menu

	if err := Database.DB.Where("menu_id = ?", menuID).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching menu",
		})
	}

	var data map[string]interface{}
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse body",
		})
	}

	allowedFields := map[string]bool{
		"name":       true,
		"category":   true,
		"start_date": true,
		"end_date":   true,
	}
	updates := make(map[string]interface{})

	for key, value := range data {
		if !allowedFields[key] {
			continue
		}
		text, ok := value.(string)
		if !ok || strings.TrimSpace(text) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": key + " must be a non-empty string",
			})
		}
		updates[key] = strings.TrimSpace(text)
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid fields to update. Allowed: name, category, start_date, end_date",
		})
	}

	if err := Database.DB.Model(&menu).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update menu",
		})
	}

	if err := Database.DB.Where("menu_id = ?", menuID).First(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Menu updated but could not be retrieved",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu updated successfully",
		"menu":    menu,
	})
}

func DeleteMenu(c fiber.Ctx) error {
	menuId := c.Params("menu_id")
	var menu Models.Menu

	// Step 1: Find menu by menu_id (NOT primary key)
	if err := Database.DB.Where("menu_id = ?", menuId).First(&menu).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Menu not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching menu",
		})
	}

	// Step 2: Delete the menu
	if err := Database.DB.Delete(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete menu",
		})
	}

	// Step 3: Return success with the deleted record's info
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Menu deleted successfully",
		"data": fiber.Map{
			"menu_id":  menu.Menu_ID,
			"name":     menu.Name,
			"category": menu.Category,
		},
	})
}
