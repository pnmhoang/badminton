package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"badminton-backend/internal/models"
	"badminton-backend/internal/views"
)

type TournamentController struct {
	db *gorm.DB
}

func NewTournamentController(db *gorm.DB) *TournamentController {
	return &TournamentController{db: db}
}

func (tc *TournamentController) GetTournaments(c *gin.Context) {
	var tournaments []models.Tournament
	if err := tc.db.Preload("Matches").Find(&tournaments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch tournaments",
		})
		return
	}

	var tournamentResponses []views.TournamentResponse
	for _, tournament := range tournaments {
		tournamentResponses = append(tournamentResponses, views.ToTournamentResponse(tournament))
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Tournaments retrieved successfully",
		Data:    tournamentResponses,
	})
}

func (tc *TournamentController) CreateTournament(c *gin.Context) {
	var tournament models.Tournament
	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	if err := tc.db.Create(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create tournament",
		})
		return
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Tournament created successfully",
		Data:    views.ToTournamentResponse(tournament),
	})
}

func (tc *TournamentController) GetTournament(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	var tournament models.Tournament
	if err := tc.db.Preload("Matches").First(&tournament, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch tournament",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Tournament retrieved successfully",
		Data:    views.ToTournamentResponse(tournament),
	})
}

func (tc *TournamentController) UpdateTournament(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	var tournament models.Tournament
	if err := tc.db.First(&tournament, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Tournament not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch tournament",
		})
		return
	}

	if err := c.ShouldBindJSON(&tournament); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	if err := tc.db.Save(&tournament).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update tournament",
		})
		return
	}

	// Load matches data
	tc.db.Preload("Matches").First(&tournament, tournament.ID)

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Tournament updated successfully",
		Data:    views.ToTournamentResponse(tournament),
	})
}

func (tc *TournamentController) DeleteTournament(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	if err := tc.db.Delete(&models.Tournament{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete tournament",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Tournament deleted successfully",
	})
}
