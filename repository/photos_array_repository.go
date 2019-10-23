package repository

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strconv"
)

type PhotosArrayRepository struct {
	dirPath string
}

func (p *PhotosArrayRepository) SavePhoto(file multipart.File, id string) (returnErr error) {
	defer func() {
		err := file.Close()

		if err != nil && returnErr == nil {
			log.Printf("An error occurred: %v", err)
			returnErr = err
		}
	}()

	tempFile, err := ioutil.TempFile(p.dirPath, "upload-*.png")
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	defer func() {
		err := tempFile.Close()

		if err != nil && returnErr == nil {
			log.Printf("An error occurred: %v", err)
			returnErr = err
		}
	}()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}
	err = os.Rename(tempFile.Name(), p.dirPath+id+".png")

	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		return err
	}

	return nil
}

func (p *PhotosArrayRepository) GetPhoto(id int) (*os.File, error) {

	fileName := strconv.Itoa(id)
	file, err := os.Open(p.dirPath + fileName + ".png")
	if err != nil {
		log.Printf("An error occurred: %v", err)
		file, err = os.Open(p.dirPath + "default" + ".png")
		return file, err
	}
	return file, nil
}

func NewPhotosArrayRepository(path string) PhotoRepository {
	return &PhotosArrayRepository{dirPath: path}
}
