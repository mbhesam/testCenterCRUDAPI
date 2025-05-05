package main

import (
	"net/http"
	"testCenterApi/db"
	"testCenterApi/docs" // import generated docs
	"testCenterApi/routes"

	"github.com/gin-contrib/cors" // CORS middleware
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // swagger middleware
)

// @title Test Center API
// @version 1.0
// @description This is a simple API for creating crud transactions..
// @host localhost:8123
// @BasePath /

func main() {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Test Center API"
	docs.SwaggerInfo.Description = "This is a simple API for creating crud transactions."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "172.23.95.0:8123"
	docs.SwaggerInfo.BasePath = "/"

	db.ConnectDatabase()

	router := gin.Default()
	router.Use(cors.Default()) // Use default CORS settings
	// Serve Swagger UI at /docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})
	// Serve your API route "/v1/table/create"
	router.POST("/v1/table/create", routes.CreateCategoryTable)
	// Serve your API route "/v1/table/read"
	router.GET("/v1/table/read", routes.ReadCategoryTable)
	// Serve your API route "/v1/table/create"
	router.POST("/v1/table/delete", routes.DeleteCategoryTable)
	// Serve your API route "/v1/table/:category/question/create"
	router.POST("/v1/table/:category/question/create", routes.CreateCategoryQuestion)
	// Serve your API route "/v1/table/:category/question/read"
	router.POST("/v1/table/:category/question/read", routes.ReadCategoryQuestion)
	// Serve your API route "/v1/table/:category/question/update"
	router.POST("/v1/table/:category/question/update", routes.UpdateCategoryQuestion)
	// Serve your API route "/v1/table/:category/question/delete"
	router.POST("/v1/table/:category/question/delete", routes.DeleteCategoryQuestion)

	// Start the server on port 8123
	router.Run(":8123")
}
