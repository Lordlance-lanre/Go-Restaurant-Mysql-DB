package Controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Database"
	utils "github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Helpers"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func validateEmail(email string) bool {
	regisMail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regisMail.MatchString(email)
}

func Signup(c fiber.Ctx) error {
	var data map[string]interface{}

	var userData Models.User

	if err := c.Bind().Body(&data); err != nil {
		fmt.Println("Unable to parse body")
	}

	//check if password is less than 6 characters

	if len(data["password"].(string)) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be at least 6 characters or more",
		})
	}
	if !validateEmail(strings.TrimSpace(data["email"].(string))) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email format",
		})
	}

	Database.DB.Where("email = ?", strings.TrimSpace(data["email"].(string))).First(&userData)

	if userData.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	userData = Models.User{
		FirstName: data["first_name"].(string),
		LastName:  data["last_name"].(string),
		Email:     strings.TrimSpace(data["email"].(string)),
		Phone:     data["phone"].(string),
		// Password: data["password"].(string),
	}

	userData.SetPassword(data["password"].(string))
	err := Database.DB.Create(&userData).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to register user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    userData,
	})
}

func Login(c fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse body",
		})
	}
	var userData Models.User
	Database.DB.Where("email = ?", strings.TrimSpace(data["email"].(string))).First(&userData)
	if userData.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email address",
		})
	}
	if err := userData.ComparePassword(data["password"].(string)); !err {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}
	token, err := utils.GenerateJWT(strconv.Itoa(int(userData.ID)))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(100 * time.Second),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged in successfully",
		"user":    userData,
		"token":   token,
	})
}

func GetAllUsers(c fiber.Ctx) error {

	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage", "10"))
	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * recordsPerPage

	var users []Models.User
	var totalUsers int64

	if err := Database.DB.Model(&Models.User{}).Count(&totalUsers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to count users",
		})
	}

	if err := Database.DB.Limit(recordsPerPage).Offset(offset).Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve users",
		})
	}

	totalPages := (int(totalUsers) + recordsPerPage - 1) / recordsPerPage

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"users":   users,
		"meta": fiber.Map{
			"total":          totalUsers,
			"page":           page,
			"recordsPerPage": recordsPerPage,
			"totalPages":     totalPages,
		},
	})
}

func GetUserByID(c fiber.Ctx) error {
	userId := c.Params("id")
	var user Models.User

	err := Database.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func UpdateUser(c fiber.Ctx) error {
	userId := c.Params("id")
	var user Models.User

	if err := Database.DB.First(&user, userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error occurred while fetching user",
		})
	}

	var data map[string]interface{}
	if err := c.Bind().Body(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unable to parse body",
		})
	}

	if firstName, ok := data["first_name"].(string); ok && strings.TrimSpace(firstName) != "" {
		user.FirstName = firstName
	}
	if lastName, ok := data["last_name"].(string); ok && strings.TrimSpace(lastName) != "" {
		user.LastName = lastName
	}
	if phone, ok := data["phone"].(string); ok && strings.TrimSpace(phone) != "" {
		user.Phone = phone
	}
	if email, ok := data["email"].(string); ok && strings.TrimSpace(email) != "" {
		email = strings.TrimSpace(email)
		if !validateEmail(email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid email format",
			})
		}

		var existingUser Models.User
		if err := Database.DB.Where("email = ? AND id != ?", email, userId).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Email already exists",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to validate email",
			})
		}

		user.Email = email
	}
	if password, ok := data["password"].(string); ok && strings.TrimSpace(password) != "" {
		if len(password) < 6 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Password must be at least 6 characters or more",
			})
		}
		user.SetPassword(password)
	}

	if err := Database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

func DeleteUser(c fiber.Ctx) error {
	fmt.Println("DeleteUser")
	return nil
}
