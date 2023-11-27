package utils

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

const projectDirName = "nat-dev"

func ConvertStringToArray(input string) []string {
	input = strings.Trim(input, `"`)

	result := strings.Split(input, `","`)

	return result
}

func LoadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentDir, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentDir))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal(err.Error())
	}
}
