package repository

import (
	"mime/multipart"
	"os"
)

type PhotoRepository interface {
	SavePhoto(file multipart.File, id string) (returnErr error)
	GetPhoto(id int) (*os.File, error)
}
