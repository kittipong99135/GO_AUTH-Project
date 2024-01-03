package controllers

import (
	"context"
	"fmt"
	"god-dev/database"
	"god-dev/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Functions Userparams after login - User routher.
func UserParams(c *fiber.Ctx) error { // Routes Post -> http://127.0.0.1:3000/api/user/params/dashboard
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Print(claims["uid"])
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Read user success",
		"id":      claims["uid"],
	})
}

// Functions List all users - User routher .
func UserList(c *fiber.Ctx) error { // Routes Get -> http://127.0.0.1:3000/api/user/

	// Connect database
	db := database.DBConn

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])

	// Recive find all users.
	var listUser []models.User
	resultUser := db.Find(&listUser, "id != ?", uid)
	if resultUser.RowsAffected == 0 { // Case : Table user is empty.
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	// Return Status200, json data
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find all user success",
		"result":  listUser,
	})
}

// Functions Read all users - User routher.
func UserRead(c *fiber.Ctx) error { // Routes Get -> http://127.0.0.1:3000/api/user/:id
	db := database.DBConn

	id := c.Params("id")

	var readUser models.User

	resultUser := db.Find(&readUser, "id = ?", id)
	if resultUser.RowsAffected == 0 { // Case : Table user is empty.
		return c.Status(503).JSON(fiber.Map{
			"status":  "warning",
			"message": "Warning : Can't find user",
			"result":  "No result.",
		})
	}

	// Return Status200, json data
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Find user success",
		"result":  readUser,
	})
}

// Functions Update all users - User routher.
func UserUpdate(c *fiber.Ctx) error { // Routes Put -> http://127.0.0.1:3000/api/user/:id
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

// Functions Remove all users - User routher.
func UserRemove(c *fiber.Ctx) error { // Routes Delete -> http://127.0.0.1:3000/api/user/:id
	db := database.DBConn
	id := c.Params("id")

	var user models.User

	db.Delete(&user, id)

	// Return Status200, json data
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Success : Remove user success",
	})
}

// Functions Active all users - User routher.
func UserActive(c *fiber.Ctx) error { // Routes Put -> http://127.0.0.1:3000/api/user/active/:id
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

// Functions Delete all users - User routher.
func UserLogout(c *fiber.Ctx) error { // Routes Post -> http://127.0.0.1:3000/api/user/:id
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	uid := fmt.Sprintf("%v", claims["uid"])
	result, err := DeleteFromRedis("access_token:" + uid)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Delete access token from redis server.",
			"error":   err.Error(),
		})
	}
	fmt.Println("Delete access toeken : " + result + " | Success")

	result, err = DeleteFromRedis("refresh_token:" + uid)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Delete refresh  token from redis server.",
			"error":   err.Error(),
		})
	}
	fmt.Println("Delete access toeken : " + result + " | Success")

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User Logout success.",
	})

}

func MiddleWareEndpoint(c *fiber.Ctx) error {
	return c.SendString("Hello, Test middleware endpoint.")
}

// Functions delete from redis.
func DeleteFromRedis(key string) (string, error) {
	rd := database.RDConn
	ctx := context.Background()
	val, err := rd.GetDel(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
