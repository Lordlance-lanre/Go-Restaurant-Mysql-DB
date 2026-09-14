package Routes

import (
	// "fmt"
	"time"

	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Controllers"
	"github.com/Lordlance-lanre/Go-Restaurant-Mysql-DB.git/Middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func AppRoutes(app *fiber.App) {
	// fmt.Println("Routes")
	// auth.Post("/login", Controllers.Login)
	// auth.Post("/logout", Middleware.AuthGuard, Controllers.Logout)

	auth := app.Group("/restaurants/api/auth")
	auth.Post("/signup", Controllers.Signup)
	loginLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many login attempts. Please try again after a minute.",
			})
		},
	})
	auth.Post("/login", loginLimiter, Controllers.Login)
	// Protect user listing with AuthGuard middleware
	protected := app.Group("/restaurants/api", Middleware.AuthGuard)
	protected.Get("/users", Controllers.GetAllUsers)
	protected.Get("/user/:id", Controllers.GetUserByID)
	protected.Put("/user/:id", Controllers.UpdateUser)

	//menu routes
	protected.Post("/create-menu", Controllers.CreateMenu)
	protected.Get("/all-menus", Controllers.GetAllMenus)
	protected.Get("/menu/:menu_id", Controllers.GetMenuByID)
	protected.Put("/menu/update/:menu_id", Controllers.UpdateMenu)
	protected.Delete("/menu/:menu_id", Controllers.DeleteMenu)

	// food routes
	protected.Post("/create-food", Controllers.CreateFoods)
	protected.Get("/all-foods", Controllers.GetAllFoods)
	protected.Get("/food/:food_id", Controllers.GetFoodByID)
	protected.Put("/food/update/:food_id", Controllers.UpdateFood)
	protected.Delete("/food/:food_id", Controllers.DeleteFood)

}
