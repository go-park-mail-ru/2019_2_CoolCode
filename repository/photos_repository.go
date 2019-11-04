package repository

import (
	"mime/multipart"
	"os"
)

//go:generate moq -out photo_repository_mock.go . PhotoRepository
type PhotoRepository interface {
	SavePhoto(file multipart.File, id string) (returnErr error)
	GetPhoto(id int) (*os.File, error)
}
