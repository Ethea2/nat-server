package services

import (
	"github.com/go-playground/validator"

	"github.com/Ethea2/nat-dev/models"
	"github.com/Ethea2/nat-dev/utils"
)

var validate = validator.New()

type mediaUpload interface {
	FileUpload(file models.File) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file models.File) (string, error) {
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	uploadUrl, err := utils.UploadPhoto(file.File)
	if err != nil {
		return "", err
	}

	return uploadUrl, nil
}
