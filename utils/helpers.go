package utils

import (
	"log"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}
