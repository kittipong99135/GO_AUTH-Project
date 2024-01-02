package controllers

import (
	"god-dev/database"
	"god-dev/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func User(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Read user success",
		"id":      claims["uid"],
	})
}

func UserList(c *fiber.Ctx) error {
	db := database.DBConn

	var listUser []models.User
	resultUser := db.Find(&listUser)
	if resultUser.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	return c.Status(500).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find all user success",
		"result":  listUser,
	})
}

func UserRead(c *fiber.Ctx) error {
	db := database.DBConn

	id := c.Params("id")

	var readUser models.User

	resultUser := db.Find(&readUser, "id = ?", id)
	if resultUser.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find user success",
		"result":  readUser,
	})
}

func UserUpdate(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Update body parser error.",
			"error":   err.Error(),
		})
	}

	updateUser := models.User{
		Email: user.Email,
		Name:  user.Name,
		Phone: user.Phone,
		Age:   user.Age,
		Rank:  user.Rank,
	}

	db.Where("id = ?", id).Updates(&updateUser)
	return c.Status(503).JSON(fiber.Map{
		"status":  "error",
		"message": "Error : Update body parser error.",
		"result":  updateUser,
	})
}

func UserRemove(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")

	var user models.User

	db.Delete(&user, id)
	return c.Status(503).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Remove user success",
	})

}

func UserActive(c *fiber.Ctx) error {
	db := database.DBConn

	id := c.Params("id")

	var user models.User

	result := db.First(&user, "id = ?", id)
	if result.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	activeUser := models.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Phone:    user.Phone,
		Age:      user.Age,
		Rank:     user.Rank,
		Status:   "active",
		Role:     user.Role,
	}

	db.Where("id = ?", id).Updates(&activeUser)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Active user success.",
		"result":  activeUser,
	})
}
