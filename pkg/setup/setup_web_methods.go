package setup

import (
	"net/http"

	"github.com/avnes/test-persistence/pkg/database"
	"github.com/avnes/test-persistence/pkg/filesystem"
	"github.com/avnes/test-persistence/pkg/logging"
	"github.com/gin-gonic/gin"
)

// Setup godoc
// @BasePath /persistence/api/v1
// @Summary Setup filesystem and database
// @Schemes
// @Description Create a table in the database and a directory on the filesystem
// @Tags persistent
// @Accept json
// @Produce json
// @Success 201 {string} OK
// @Router /setup [post]
func Setup(c *gin.Context) {
	log := logging.SugarLogger()
	err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	errDir := filesystem.CreateDirectory()
	if errDir != nil {
		log.Fatal(errDir)
	}

	c.IndentedJSON(http.StatusCreated, "OK")
}
