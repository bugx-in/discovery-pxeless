package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"fpi/api"
)

func main() {
	// Validate required environment variables.
	requiredVariables := []string{"DISCOVERY_REMASTER", "DISCOVERY_BASE_IMAGE", "IMAGES_PATH"}

	for _, req := range requiredVariables {
		if os.Getenv(req) == "" {
			log.Fatal().Msgf("Environment variable is missing: %s", req)

		}
	}

	if err := api.Run(); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
