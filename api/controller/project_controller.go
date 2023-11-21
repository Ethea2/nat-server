package controller

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Ethea2/nat-dev/database"
	"github.com/Ethea2/nat-dev/models"
	"github.com/Ethea2/nat-dev/services"
	"github.com/Ethea2/nat-dev/utils"
)

func GetProjects(c *fiber.Ctx) error {
	query := `SELECT * FROM projects`

	var projects []models.Projects
	projRes, err := database.DataBase.Query(query)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}
	for projRes.Next() {
		var tempProj models.Projects
		err := projRes.Scan(
			&tempProj.ID,
			&tempProj.UserID,
			&tempProj.Title,
			&tempProj.Body,
			&tempProj.Image,
		)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				StatusCode: 500,
				Message:    "error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			})
		}
		projects = append(projects, tempProj)
	}

	techstackQuery := `SELECT techstacks FROM techstacks WHERE projectID = (?)`

	var projectWithTstacks []models.Projects

	for _, project := range projects {
		projectID := project.ID
		var tech []string
		res, err := database.DataBase.Query(techstackQuery, projectID)
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

		project.Techstacks = tech

		projectWithTstacks = append(projectWithTstacks, project)
	}

	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "success",
		Data: &fiber.Map{
			"data": projectWithTstacks,
		},
	})
}

func CreateProject(c *fiber.Ctx) error {
	userID := c.Params("userID")

	formHeader, err := c.FormFile("image")
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": "There was an error with parsing the image!",
			},
		})
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	uploadUrl, err := services.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	title := c.FormValue("title")
	body := c.FormValue("body")
	stringTechstacks := c.FormValue("techstacks")
	techstacks := utils.ConvertStringToArray(stringTechstacks)
	image := uploadUrl
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": "There was an in the userID convertion",
			},
		})
	}

	project := models.Projects{
		UserID:     intUserID,
		Title:      title,
		Body:       body,
		Techstacks: techstacks,
		Image:      image,
	}
	query := `INSERT INTO projects (userID, title, body, image) VALUES (?, ?, ?, ?)`

	res, err := database.DataBase.Exec(
		query,
		project.UserID,
		project.Title,
		project.Body,
		project.Image,
	)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}
	projectID64, err := res.LastInsertId()
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: http.StatusInternalServerError,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}
	project.ID = int(projectID64)

	techstackQuery := `INSERT INTO techstacks (projectID, techstacks) VALUES (?, ?)`

	for _, tech := range techstacks {
		res, err := database.DataBase.Exec(techstackQuery, project.ID, tech)
		if err != nil {
			return c.Status(500).JSON(models.Response{
				StatusCode: http.StatusInternalServerError,
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
			"data": project,
		},
	})
}

func goGetTechstacks(
	c *fiber.Ctx,
	projectID int,
	techChan chan<- *[]string,
	errorChan chan<- error,
) {
	techstacks := []string{}
	techstackQuery := `SELECT techstacks FROM techstacks WHERE projectID = (?)`
	responseRows, err := database.DataBase.Query(techstackQuery, int(projectID))
	if err != nil {
		errorChan <- c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
		close(techChan)
		return
	}
	defer responseRows.Close()

	for responseRows.Next() {
		var techstack string
		err := responseRows.Scan(&techstack)
		if err != nil {
			errorChan <- c.Status(500).JSON(models.Response{
				StatusCode: 500,
				Message:    "error",
				Data: &fiber.Map{
					"data": err.Error(),
				},
			})
			close(techChan)
			return
		}
		techstacks = append(techstacks, techstack)
	}
	techChan <- &techstacks
	close(techChan)
}

func goGetProject(
	c *fiber.Ctx,
	projectID int,
	projectChan chan<- *models.Projects,
	errorChan chan<- error,
) {
	query := `SELECT * FROM projects WHERE id = (?)`
	project := new(models.Projects)

	scanErr := database.DataBase.QueryRow(query, projectID).
		Scan(&project.ID, &project.UserID, &project.Title, &project.Body, &project.Image)
	if scanErr != nil {
		errorChan <- c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": scanErr.Error(),
			},
		})
		close(projectChan)
		return
	}
	projectChan <- project
	close(projectChan)
}

func GetSinglePost(c *fiber.Ctx) error {
	stringProjectID := c.Params("id")
	projectID, err := strconv.Atoi(stringProjectID)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	errorChan := make(chan error)
	projectChan := make(chan *models.Projects)
	techChan := make(chan *[]string)

	go goGetProject(c, projectID, projectChan, errorChan)
	go goGetTechstacks(c, projectID, techChan, errorChan)

	techstacks := <-techChan
	if techstacks == nil {
		return <-errorChan
	}

	project := <-projectChan
	if project == nil {
		return <-errorChan
	}

	project.Techstacks = *techstacks

	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "success",
		Data: &fiber.Map{
			"data": project,
		},
	})
}

func UpdateProject(c *fiber.Ctx) error {
	stringProjectID := c.Params("id")
	projectID, err := strconv.Atoi(stringProjectID)
	if err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	projects := new(models.Projects)

	if err := c.BodyParser(projects); err != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	projQuery := `SELECT * FROM projects WHERE id = (?)`

	oldProj := new(models.Projects)

	reqErr := database.DataBase.QueryRow(projQuery, projectID).
		Scan(&oldProj.ID, &oldProj.UserID, &oldProj.Title, &oldProj.Body, &oldProj.Image)
	if reqErr != nil {
		return c.Status(500).JSON(models.Response{
			StatusCode: 500,
			Message:    "error",
			Data: &fiber.Map{
				"data": err.Error(),
			},
		})
	}

	test := reflect.ValueOf(projects)
	num := test.NumField()
	for i := 0; i < num; i++ {
		project := test.Field(i)
		fmt.Println(project)
	}

	return c.Status(200).JSON(models.Response{
		StatusCode: 200,
		Message:    "error",
		Data: &fiber.Map{
			"data": oldProj,
		},
	})
}
