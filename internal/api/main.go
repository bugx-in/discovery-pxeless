package api

import (
	"fpi/internal/docs"
	"fpi/internal/fpi"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"os"
)

var imagesPath string = os.Getenv("IMAGES_PATH")

// Generates a discovery image with the customized data.
// Puts the image available to download by the server.
func generateDiscoveryImage(c *gin.Context) {
	// Initialize discovery image.
	discoveryImage := fpi.NewDiscoveryImage()

	// Call BindJSON to bind the received JSON request and tranform it to a structure.
	if err := c.BindJSON(&discoveryImage); err != nil {

		// Return error if data could not be serialized.
		log.Error().Msgf("Received data could not be serialized: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Validate some requirements.
	_, error := fpi.DiscoveryImageValidate(&discoveryImage)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": error.Error()})
	}

	// Generate the discovery image.
	result := fpi.GenerateDiscoveryImage(&discoveryImage, imagesPath)

	log.Debug().Msg(result)

	c.JSON(http.StatusOK, gin.H{"result": "Image will be ready in 2 minutes."})
}

// Get API status
// @BasePath /api/v1
// List images godoc
// @Summary Get API status
// @Schemes
// @Description API status
// @Tags Status
// @Accept json
// @Produce json
// @Success 200 {dict} Status of the API.
// @Router /api/v1/status [get]
func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Online"})
}

// List the images available.
// @BasePath /api/v1
// List images godoc
// @Summary List ISO images
// @Schemes
// @Description List all ISO images available.
// @Tags Images
// @Accept json
// @Produce json
// @Success 200 {dict} Dictionary with the results.
// @Router /api/v1/images [get]
func listImages(c *gin.Context) {
	// Get all images.
	isos := fpi.ListImages(imagesPath)

	c.JSON(http.StatusOK, gin.H{"results": isos})
}

// Puts the ISO image on the stream for clients to download it.
// @BasePath /api/v1
// Get Image godoc
// @Summary Get ISO image.
// @Schemes
// @Description Get an ISO image.
// @Tags Images
// @Accept json
// @Produce json
// @Success 200 {stream} ISO image
// @Router /api/v1/images/{name} [get]
func getImage(c *gin.Context) {
	// Check if image exists.
	if !fpi.ImageExist(c.Param("name"), imagesPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image does not exist: " + c.Param("name")})
	}

	// Return the image.
	c.FileAttachment(imagesPath+"/"+c.Param("name"), c.Param("name"))
}

// Delete and image.
// @BasePath /api/v1
// @Summary Delete an image.
// @Tags Images
// @Accept json
// @Produce json
// @Success 200 {string}e
// @Router /api/v1/images/{name} [delete]
func deleteImage(c *gin.Context) {
	// Check if image exists.
	if !fpi.ImageExist(c.Param("name"), imagesPath) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image does not exist: " + c.Param("name")})
	}

	// Remove the image and reurn the result.
	os.Remove(imagesPath + "/" + c.Param("name"))
	c.JSON(http.StatusOK, gin.H{"result": "Image deleted: " + c.Param("name")})
	return
}

func Run() {
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1/"

	// V1 version.
	v1 := router.Group("/api/v1/")
	{
		v1.GET("/status", status)
		v1.GET("/images", listImages)
		v1.HEAD("/images/:name", getImage)
		v1.GET("/images/:name", getImage)
		v1.DELETE("/images/:name", deleteImage)
		v1.POST("/generate", generateDiscoveryImage)
	}
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
