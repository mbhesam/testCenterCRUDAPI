package routes

import (
	"encoding/json"
	"net/http"
	"testCenterApi/db"
	"testCenterApi/models"

	"github.com/gin-gonic/gin"
)

// CreateQuestionRequest represents a question with options and answers
type CreateQuestionRequest struct {
	Question string   `json:"question" example:"Which of the following are valid IPv4 private address ranges?"`
	Options  []string `json:"options" example:"10.0.0.0 – 10.255.255.255,172.16.0.0 – 172.31.255.255,192.168.0.0 – 192.168.255.255,169.254.0.0 – 169.254.255.255"`
	Answers  []int    `json:"answers" example:"0,1,2"`
}

type ReadQuestionRequest struct {
	Include string `json:"include" example:"private"` // Optional filter substring
}

type UpdateQuestionRequest struct {
	ID       uint     `json:"id" binding:"required" example:"1"`
	Question string   `json:"question" example:"Which of the following are valid IPv4 private address ranges?"`
	Options  []string `json:"options" example:"10.0.0.0 – 10.255.255.255,172.16.0.0 – 172.31.255.255,192.168.0.0 – 192.168.255.255,169.254.0.0 – 169.254.255.255"`
	Answers  []int    `json:"answers" example:"0,1,2"`
}

type DeleteQuestionRequest struct {
	ID uint `json:"id" binding:"required" example:"1"`
}

/////version V1 API

// CreateCategoryQuestion godoc
// @Summary Create a new question in a category
// @Description Create a question with dynamic options and answers
// @Tags questions
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param createQuestion body CreateQuestionRequest true "Question Input"
// @Param createQuestion body routes.CreateQuestionRequest true "Question Input" example({"question": "Which of the following are valid IPv4 private address ranges?", "options": ["10.0.0.0 – 10.255.255.255", "172.16.0.0 – 172.31.255.255", "192.168.0.0 – 192.168.255.255"], "answers": [0, 1, 2]})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/{category}/question/create [post]
func CreateCategoryQuestion(c *gin.Context) {
	category := c.Param("category") // e.g., ccna, ccnp
	var req CreateQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Options) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "minimum 2 options required"})
		return
	}

	for _, a := range req.Answers {
		if a < 0 || a >= len(req.Options) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "answer index out of range"})
			return
		}
	}

	optionsJSON, _ := json.Marshal(req.Options)
	answersJSON, _ := json.Marshal(req.Answers)

	question := models.DynamicQuestion{
		Question: req.Question,
		Options:  optionsJSON,
		Answers:  answersJSON,
	}

	if err := db.DB.Table(category).Create(&question).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "question saved", "category": category})
}

// ReadCategoryQuestion godoc
// @Summary Read all questions from a category table
// @Description Fetches questions from a category table, optionally filtered by substring match in question
// @Tags questions
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param request body routes.ReadQuestionRequest true "Optional Include Filter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/{category}/question/read [post]
func ReadCategoryQuestion(c *gin.Context) {
	category := c.Param("category")
	var req ReadQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rows []map[string]interface{}
	query := db.DB.Table(category)

	if req.Include != "" {
		query = query.Where("question ILIKE ?", "%"+req.Include+"%")
	}

	if err := query.Find(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rows})
}

// UpdateCategoryQuestion godoc
// @Summary Update a question in a category table
// @Description Updates question text, options, and answers by ID
// @Tags questions
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param request body routes.UpdateQuestionRequest true "Question Update Input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/{category}/question/update [post]
func UpdateCategoryQuestion(c *gin.Context) {
	category := c.Param("category")
	var req UpdateQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	optionsJSON, _ := json.Marshal(req.Options)
	answersJSON, _ := json.Marshal(req.Answers)

	updateData := map[string]interface{}{
		"question": req.Question,
		"options":  optionsJSON,
		"answers":  answersJSON,
	}

	if err := db.DB.Table(category).Where("id = ?", req.ID).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "question updated", "id": req.ID})
}

// DeleteCategoryQuestion godoc
// @Summary Delete a question by ID
// @Description Deletes a question from a specific category table
// @Tags questions
// @Accept json
// @Produce json
// @Param category path string true "Category name"
// @Param request body routes.DeleteQuestionRequest true "Delete Input"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /v1/table/{category}/question/delete [post]
func DeleteCategoryQuestion(c *gin.Context) {
	category := c.Param("category")
	var req DeleteQuestionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Table(category).Where("id = ?", req.ID).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "question deleted", "id": req.ID})
}
