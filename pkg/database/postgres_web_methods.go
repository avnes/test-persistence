package database

import (
	"net/http"

	"github.com/avnes/test-persistence/pkg/logging"
	"github.com/gin-gonic/gin"
)

type configuration struct {
	NumberOfInserts int `json:"number_of_inserts" form:"number_of_inserts" binding:"required" default:"10"`
}

// GetRandomData godoc
// @BasePath /persistence/api/v1
// @Summary Get random data
// @Schemes
// @Description Return a list of random data from a database. Random data is a GUID and a timestamp
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {object} database.randomData
// @Router /database [get]
func GetRandomData(c *gin.Context) {
	log := logging.SugarLogger()
	data, err := httpGetRandomData()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, data)
}

// PostRandomData godoc
// @BasePath /persistence/api/v1
// @Summary Post random data
// @Schemes
// @Description Inserts new records into the database
// @Tags persistent
// @Accept multipart/form-data
// @Produce json
// @Param   configuration	formData	configuration  true  "How many rows would you like to insert?"
// @Success 201 {boolean} true
// @Router /database [post]
func PostRandomData(c *gin.Context) {
	log := logging.SugarLogger()
	retVal := true
	var config configuration
	err := c.Bind(&config)
	if err != nil {
		log.Fatal(err)
		retVal = false
	}
	popErr := populateTable(config.NumberOfInserts)
	if popErr != nil {
		log.Fatal(popErr)
		retVal = false
	}
	c.IndentedJSON(http.StatusCreated, retVal)
}

// DeleteRandomData godoc
// @BasePath /persistence/api/v1
// @Summary Delete random data
// @Schemes
// @Description Delete all the random test data
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {boolean} true
// @Router /database [delete]
func DeleteRandomData(c *gin.Context) {
	log := logging.SugarLogger()
	retVal := true
	err := truncateTable()
	if err != nil {
		log.Fatal(err)
		retVal = false
	}
	c.IndentedJSON(http.StatusOK, retVal)
}

// GetRandomDataCount godoc
// @BasePath /persistence/api/v1
// @Summary Get random data count
// @Schemes
// @Description Return a count of rows in the random data table
// @Tags persistent
// @Accept json
// @Produce json
// @Success 200 {object} common.Counter
// @Router /database/count [get]
func GetRandomDataCount(c *gin.Context) {
	log := logging.SugarLogger()
	data, err := httpGetRandomDataCount()
	if err != nil {
		log.Fatal(err)
	}
	c.IndentedJSON(http.StatusOK, data)
}
