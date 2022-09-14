package filesystem

import (
	"log"
	"net/http"
	"os"

	"github.com/avnes/test-persistence/pkg/database"
	"github.com/avnes/test-persistence/pkg/logging"
	"github.com/gin-gonic/gin"
)

type configuration struct {
	SizeInBytes   int64 `json:"size_in_bytes" form:"size_in_bytes" binding:"required" default:"10485760"`
	NumberOfFiles int   `json:"number_of_files" form:"number_of_files" binding:"required" default:"5"`
}

// GetFiles godoc
// @BasePath /persistence/api/v1
// @Summary Get files
// @Schemes
// @Description Return a list of files from the filesystem with name, size, modification date and file type information.
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {object} filesystem.metadata
// @Router /files [get]
func GetFiles(c *gin.Context) {
	log := logging.SugarLogger()
	files, err := httpGetFiles()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, files)
}

// PostFiles godoc
// @BasePath /persistence/api/v1
// @Summary Post files
// @Schemes
// @Description Create new files on the filesystem
// @Tags persistent
// @Accept multipart/form-data
// @Produce json
// @Param   configuration	formData	configuration  true  "Where are the files located?"
// @Success 201 {object} filesystem.metadata
// @Router /files [post]
func PostFiles(c *gin.Context) {
	log := logging.SugarLogger()
	var config configuration
	err := c.Bind(&config)
	if err != nil {
		log.Fatal(err)
	}
	files, err := httpPostFiles(config.NumberOfFiles, config.SizeInBytes)
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusCreated, files)
}

// DeleteFiles godoc
// @BasePath /persistence/api/v1
// @Summary Delete files
// @Schemes
// @Description Delete the test files
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {boolean} true
// @Router /files [delete]
func DeleteFiles(c *gin.Context) {
	log := logging.SugarLogger()
	files, err := httpDeleteFiles()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, files)
}

// GetFilesCount godoc
// @Summary Get files count
// @Schemes
// @Description Return a count of the number of files
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {object} common.Counter
// @Router /files/count [get]
func GetFilesCount(c *gin.Context) {
	log := logging.SugarLogger()
	files, err := httpGetFilesCount()
	if err != nil {
		log.Fatal(err)

	}
	c.IndentedJSON(http.StatusOK, files)
}

func GetIndex(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/persistence/swagger/index.html")
}

func GetHealth(c *gin.Context) {
	_, err := os.Stat(getDirectory())
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
	db, err := database.GetConnection()
	_ = db
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"message": "OK",
	})
}
