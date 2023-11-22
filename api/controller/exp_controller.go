package controller

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Ethea2/nat-dev/database"
	"github.com/Ethea2/nat-dev/models"
)

func CreateExperience(c *fiber.Ctx) error {
	stringUserID := c.Params("id")
	exp := new(models.Experiences)
	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		return c.Status(200).JSON(models.Response{
			StatusCode: 200,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	exp.UserID = userID

	if err := c.BodyParser(&exp); err != nil {
		return c.Status(200).JSON(models.Response{
			StatusCode: 200,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	fmt.Print(exp)

	query := `INSERT INTO experiences (userID, title, body, position) VALUES (?, ?, ?, ?)`
	res, err := database.DataBase.Exec(query, exp.UserID, exp.Title, exp.Body, exp.Position)
	if err != nil {
		return c.Status(200).JSON(models.Response{
			StatusCode: 200,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	newExpID, err := res.LastInsertId()
	if err != nil {
		return c.Status(200).JSON(models.Response{
			StatusCode: 200,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	techstacksQuery := `INSERT INTO techstacks (expID, techstacks) VALUES (?, ?)`

	for _, tech := range exp.Techstacks {
		res, err := database.DataBase.Exec(techstacksQuery, newExpID, tech)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				StatusCode: 500,
				Message:    "error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			})
		}
		fmt.Print(res.LastInsertId())
	}
	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "success",
		Data: &fiber.Map{
			"data": newExpID,
		},
	})
}

func GetExperiences(c *fiber.Ctx) error {
	query := `SELECT * FROM experiences`

	var experiences []models.Experiences

	res, err := database.DataBase.Query(query)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	for res.Next() {
		var temp models.Experiences
		err := res.Scan(&temp.ID, &temp.UserID, &temp.Title, &temp.Body, &temp.Position)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				StatusCode: 500,
				Message:    "error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			})
		}
		experiences = append(experiences, temp)
	}

	var expWithTechstacks []models.Experiences
	techstackQuery := `SELECT techstacks FROM techstacks WHERE expID = (?)`
	for _, exp := range experiences {
		expID := exp.ID

		var tech []string
		res, err := database.DataBase.Query(techstackQuery, expID)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				StatusCode: 500,
				Message:    "error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			})
		}
		for res.Next() {
			var temp string
			err := res.Scan(&temp)
			if err != nil {
				return c.Status(500).JSON(models.Response{
					StatusCode: 500,
					Message:    "error",
					Data: &fiber.Map{
						"data": err.Error(),
					},
				})
			}
			tech = append(tech, temp)
		}
		exp.Techstacks = tech

		expWithTechstacks = append(expWithTechstacks, exp)
	}

	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "success",
		Data: &fiber.Map{
			"data": expWithTechstacks,
		},
	})
}
