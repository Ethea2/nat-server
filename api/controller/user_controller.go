package controller

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ethea2/nat-dev/database"
	"github.com/Ethea2/nat-dev/models"
)

func Login(c *fiber.Ctx) error {
	godotenv.Load(".env")
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	login := new(LoginInput)
	if err := c.BodyParser(login); err != nil {
		return c.SendString(err.Error())
	}

	username := login.Username
	// password := login.Password
	//
	user := new(models.User)
	query := `SELECT * FROM user WHERE username = (?)`

	err := database.DataBase.QueryRow(query, username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return c.Status(400).JSON(models.Response{
			StatusCode: 400,
			Message:    "error",
			Data: &fiber.Map{
				"data": "No body!",
			},
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		return c.Status(400).JSON(models.Response{
			StatusCode: 400,
			Message:    "error",
			Data: &fiber.Map{
				"data": "Wrong password!",
			},
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = username
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	finalToken, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "success",
		Data: &fiber.Map{
			"data": &fiber.Map{
				"token":    finalToken,
				"username": user.Username,
				"id":       user.ID,
			},
		},
	})
}

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return c.SendString(err.Error())
	}

	query := `INSERT INTO user (username, password) VALUES (?, ?)`

	res, err := database.DataBase.Exec(query, user.Username, hashedPassword)
	if err != nil {
		return c.SendString(err.Error())
	}

	fmt.Print(res.LastInsertId())

	return c.SendStatus(200)
}
