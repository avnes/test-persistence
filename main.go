package main

import (
	"github.com/avnes/test-persistence/docs"
	"github.com/avnes/test-persistence/pkg/database"
	"github.com/avnes/test-persistence/pkg/filesystem"
	"github.com/avnes/test-persistence/pkg/setup"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	docs.SwaggerInfo.BasePath = "/persistence/api/v1"
	v1 := router.Group("/persistence/api/v1")
	{
		v1.POST("/setup", setup.Setup)
		v1.GET("/files", filesystem.GetFiles)
		v1.POST("/files", filesystem.PostFiles)
		v1.DELETE("/files", filesystem.DeleteFiles)
		v1.GET("/files/count", filesystem.GetFilesCount)
		v1.GET("/database", database.GetRandomData)
		v1.POST("/database", database.PostRandomData)
		v1.DELETE("/database", database.DeleteRandomData)
		v1.GET("/database/count", database.GetRandomDataCount)
	}
	router.GET("/persistence/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/persistence/healthz", filesystem.GetHealth)
	router.GET("/persistence/", filesystem.GetIndex)
	router.Run("0.0.0.0:8080")

}
