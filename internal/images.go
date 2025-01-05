package fpi

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"slices"
	"strings"
)

type Images struct {
	imagesPath string
}

func NewImages(imagesPath string) *Images {
	return &Images{
		imagesPath: imagesPath,
	}
}

func verifyImageName(image string) error {
	if strings.Contains(image, "..") || strings.Contains(image, "/") {
		log.Warn().Msg("Request wanted to scape the folder.")
		return errors.New("Bad image name.")
	}
	return nil
}

func (i *Images) ImageExist(image string) bool {
	// Check that they don't want to escape the current folder.
	if err := verifyImageName(image); err != nil {
		return false
	}

	return slices.Contains(i.ListImages(), image)
}

// Get a list of available images.
func (i *Images) ListImages() []string {
	var isos []string

	// Read directory.
	files, err := os.ReadDir(i.imagesPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		return isos
	}

	// Add isos to list.
	for _, file := range files {
		if !file.IsDir() {
			isos = append(isos, file.Name())
		}
	}

	return isos
}

func (i *Images) DeleteImage(image string) error {
	// Check that they don't want to escape the current folder.
	if err := verifyImageName(image); err != nil {
		return err
	}
	// Remove the image and reurn the result.
	if err := os.Remove(fmt.Sprintf("%s/%s", i.imagesPath, image)); err != nil {
		return err
	}
	return nil
}
