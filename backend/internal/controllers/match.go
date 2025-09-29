package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"badminton-backend/internal/models"
	"badminton-backend/internal/views"
)

type MatchController struct {
	db *gorm.DB
}

func NewMatchController(db *gorm.DB) *MatchController {
	return &MatchController{db: db}
}

func (mc *MatchController) GetMatches(c *gin.Context) {
	var matches []models.Match
	if err := mc.db.Preload("Player1").Preload("Player2").Preload("Player3").Preload("Player4").Preload("Tournament").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch matches",
		})
		return
	}

	var matchResponses []views.MatchResponse
	for _, match := range matches {
		matchResponses = append(matchResponses, views.ToMatchResponse(match))
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Matches retrieved successfully",
		Data:    matchResponses,
	})
}

func (mc *MatchController) CreateMatch(c *gin.Context) {
	var match models.Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Set default match date if not provided
	if match.MatchDate.IsZero() {
		match.MatchDate = time.Now()
	}

	if err := mc.db.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create match",
		})
		return
	}

	// Load related data
	mc.db.Preload("Player1").Preload("Player2").Preload("Player3").Preload("Player4").Preload("Tournament").First(&match, match.ID)

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Match created successfully",
		Data:    views.ToMatchResponse(match),
	})
}

func (mc *MatchController) GetMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid match ID",
		})
		return
	}

	var match models.Match
	if err := mc.db.Preload("Player1").Preload("Player2").Preload("Player3").Preload("Player4").Preload("Tournament").First(&match, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Match not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch match",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Match retrieved successfully",
		Data:    views.ToMatchResponse(match),
	})
}

func (mc *MatchController) UpdateMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid match ID",
		})
		return
	}

	var match models.Match
	if err := mc.db.First(&match, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Match not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch match",
		})
		return
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	if err := mc.db.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update match",
		})
		return
	}

	// Load related data
	mc.db.Preload("Player1").Preload("Player2").Preload("Player3").Preload("Player4").Preload("Tournament").First(&match, match.ID)

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Match updated successfully",
		Data:    views.ToMatchResponse(match),
	})
}

func (mc *MatchController) DeleteMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid match ID",
		})
		return
	}

	if err := mc.db.Delete(&models.Match{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete match",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Match deleted successfully",
	})
}
