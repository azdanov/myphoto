package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Image is stored on local filesystem
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	imgURL := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return imgURL.String()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	Create(galleryID uint, src io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	DeleteGallery(galleryID uint) error
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, src io.Reader, filename string) error {
	path, err := is.mkdirGallery(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	imgPaths, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	images := make([]Image, len(imgPaths))
	for i := range imgPaths {
		images[i] = Image{
			Filename:  strings.Replace(imgPaths[i], path, "", 1),
			GalleryID: galleryID,
		}
	}
	return images, nil
}

func (is *imageService) DeleteGallery(galleryID uint) error {
	return os.RemoveAll(is.imagePath(galleryID))
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) mkdirGallery(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0o755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
