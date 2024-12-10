package fpi

import (
	"strings"
	"slices"
	"os"
	"github.com/rs/zerolog/log"
)

func ImageExist(image string, imagesPath string) bool {
	// Check that they don't want to escape the current folder.
	if strings.Contains(image, "..") || strings.Contains(image, "/") {
		log.Warn().Msg("Request wanted to scape the folder.")
		return false
	}

	isos := ListImages(imagesPath)

	return slices.Contains(isos, image)
}

// Get a list of available images.
func ListImages(imagesPath string) []string {
	var isos []string

	files, err := os.ReadDir(imagesPath)

	if err != nil {
		log.Warn().Msg(err.Error())
		return isos
    }

    for _, file := range files {
		if !file.IsDir() {
            isos = append(isos, file.Name())
		}
    }
	return isos
}