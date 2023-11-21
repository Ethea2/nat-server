package utils

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

func UploadPhoto(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	godotenv.Load()

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUD_NAME"),
		os.Getenv("CLOUD_KEY"),
		os.Getenv("CLOUD_SECRET"),
	)
	if err != nil {
		log.Fatal("Failed to connect to cloudinary!", err)
		return "", err
	}

	uploadedFile, err := cld.Upload.Upload(
		ctx,
		input,
		uploader.UploadParams{Folder: "Portfolio-Tests"},
	)

	return uploadedFile.SecureURL, nil
}
