package controllers

import (
	"context"
	"god-dev/database"
	"god-dev/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Regis(c *fiber.Ctx) error {
	db := database.DBConn
	var regisBody models.User

	err := c.BodyParser(&regisBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Regid body invalid.",
			"error":   err.Error(),
		})
	}

	var userExists models.User
	result := db.Find(&userExists, "email = ?", strings.TrimSpace(regisBody.Email))
	if result.RowsAffected != 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Email exists.",
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(regisBody.Password), 10)
	if err != nil {
		c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Invalid password hashing.",
			"error":   err.Error(),
		})
	}

	userRegisted := models.User{
		Email:    regisBody.Email,
		Password: "secretpass:" + string(hash),
		Name:     regisBody.Name,
		Phone:    regisBody.Phone,
		Age:      regisBody.Age,
		Rank:     regisBody.Rank,
		Status:   "nactive",
		Role:     "user",
	}

	db.Create(&userRegisted)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"massage": "User register success",
		"detail":  userRegisted,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DBConn

	var loginBody models.RequestLogin
	err := c.BodyParser(&loginBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Logging body invalid.",
			"error":   err.Error(),
		})
	}

	var user models.User
	result := db.Find(&user, "email = ?", strings.TrimSpace(loginBody.Email))
	if result.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Email invalid.",
		})
	}

	splitPass := strings.Split(user.Password, ":")[1:][0]
	err = bcrypt.CompareHashAndPassword([]byte(splitPass), []byte(loginBody.Password))
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Compare password error.",
			"error":   err.Error(),
		})
	}

	udid := strconv.Itoa(int(user.ID))

	acc_token, err := CreateToken(udid, "JWT_SECRET")
	if err != nil {
		c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Create access token  error.",
			"error":   err.Error(),
		})
	}
	SetAccessToken("access_token:"+udid, acc_token)

	rfh_token, err := CreateToken(udid, "JWT_REFRESH")
	if err != nil {
		c.Status(503).JSON(fiber.Map{
			"status":  "error",
			"message": "Error : Create refresh token error.",
			"error":   err.Error(),
		})
	}
	SerRefreshToken("refresh_token:"+udid, rfh_token)

	return c.Status(200).JSON(fiber.Map{
		"status":        "success",
		"message":       "Success : Logging in success.",
		"user":          udid,
		"token_access":  GetToken("access_token:" + udid),
		"token_refresh": GetToken("refresh_token:" + udid),
	})
}

func CreateToken(udid string, env string) (string, error) {
	cliams := jwt.MapClaims{"uid": udid}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString([]byte(os.Getenv("env")))
}

func SetAccessToken(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, time.Hour*2)
}

func GetToken(key string) string {
	rd := database.RDConn
	ctx := context.Background()
	val, _ := rd.Get(ctx, key).Result()
	return val
}

func SerRefreshToken(key string, token string) {
	rd := database.RDConn
	ctx := context.Background()
	rd.Set(ctx, key, token, 0)
}
