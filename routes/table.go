package routes

import (
	"fmt"
	"net/http"
	"testCenterApi/db"
	"testCenterApi/models"

	"github.com/gin-gonic/gin"
)

type CreateTableRequest struct {
	Category string `json:"category" binding:"required,alphanum"  example:"ccnp"` // like "ccna"
}

// ReadCategoryTableRequest defines the request structure
type TableInfo struct {
	TableName string `gorm:"column:tablename"`
}

var rawResults []TableInfo

// DeleteCategoryTableRequest defines the request structure
type DeleteCategoryTableRequest struct {
	Category string `json:"category" binding:"required,alphanum"  example:"ccnp"`
}

/////version V1 API

// CreateCategoryTable godoc
// @Summary Create a new category table
// @Description Create a new category as a dynamic table
// @Tags categories
// @Accept json
// @Produce json
// @Param request body routes.CreateTableRequest true "Category Input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/create [post]
func CreateCategoryTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new table with the given category name
	err := db.DB.Table(req.Category).AutoMigrate(&models.DynamicQuestion{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table created", "table": req.Category})
}

// ReadCategoryTable godoc
// @Summary List all question category tables
// @Description Retrieve all user-created question category tables
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]string
// @Router /v1/table/read [get]
func ReadCategoryTable(c *gin.Context) {
	type TableInfo struct {
		TableName string `gorm:"column:tablename"`
	}

	var rawResults []TableInfo

	if err := db.DB.Raw(`
		SELECT tablename 
		FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT LIKE 'schema_migrations'
	`).Scan(&rawResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var tableNames []string
	for _, t := range rawResults {
		tableNames = append(tableNames, t.TableName)
	}

	c.JSON(http.StatusOK, gin.H{"tables": tableNames})
}

// DeleteCategoryTable godoc
// @Summary Delete a category table
// @Description Drops a table belonging to a specific category (like "ccna", "ccnp")
// @Tags categories
// @Accept json
// @Produce json
// @Param request body routes.DeleteCategoryTableRequest true "Category Input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/delete [post]
func DeleteCategoryTable(c *gin.Context) {
	var req DeleteCategoryTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sql := fmt.Sprintf("DROP TABLE IF EXISTS \"%s\";", req.Category)
	if err := db.DB.Exec(sql).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table deleted", "table": req.Category})
}
