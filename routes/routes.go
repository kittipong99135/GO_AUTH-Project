package routes

import (
	"god-dev/controllers"
	"god-dev/middleware"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", controllers.Regis)
	auth.Post("login", controllers.Login)

	user := api.Group("/user", middleware.RequestAuth(), middleware.RefreshAuth())
	user.Get("/", controllers.UserList)
	user.Get("/:id", controllers.UserRead)
	user.Put("/:id", controllers.UserUpdate)
	user.Delete("/:id", controllers.UserRemove)

	user.Get("/dashboard", controllers.User)
	user.Put("/active/:id", controllers.UserActive)

}
