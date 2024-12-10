package main

import (
	"os"
	"fpi/internal/api"
	"github.com/rs/zerolog/log"
)

func main() {
	// Validate required environment variables.
	requiredVariables := []string{"DISCOVERY_REMASTER", "DISCOVERY_BASE_IMAGE", "IMAGES_PATH"}

    for _, req := range requiredVariables {
		if os.Getenv(req) == "" {
			log.Fatal().Msgf("Environment variable is missing: %s", req)

			}
	}

	api.Run()

	return 
}