package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"badminton-backend/internal/models"
	"badminton-backend/internal/views"
)

type PlayerController struct {
	db *gorm.DB
}

func NewPlayerController(db *gorm.DB) *PlayerController {
	return &PlayerController{db: db}
}

// GetPlayers now returns User objects with role="player"
// Legacy Player Controller - now works with User model
// This controller maintains backward compatibility for existing API endpoints
// but uses the new User system under the hood

func (pc *PlayerController) GetPlayers(c *gin.Context) {
	var users []models.User
	// Only get users with player role
	if err := pc.db.Where("role = ?", "player").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch players",
		})
		return
	}

	var playerResponses []views.PlayerResponse
	for _, user := range users {
		// Convert User to PlayerResponse for backward compatibility
		playerResponse := views.PlayerResponse{
			ID:      user.ID,
			Name:    user.FullName,
			Email:   user.Email,
			Ranking: user.Ranking,
		}
		playerResponses = append(playerResponses, playerResponse)
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Players retrieved successfully",
		Data:    playerResponses,
	})
}

// CreatePlayer creates a new User with player role (legacy endpoint)
func (pc *PlayerController) CreatePlayer(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
		Ranking  int    `json:"ranking"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Create User with player role
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.Name,
		Role:     models.RolePlayer,
		Ranking:  req.Ranking,
		IsActive: true,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "password_hash_failed",
			Message: "Failed to hash password",
		})
		return
	}

	if err := pc.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create player",
		})
		return
	}

	playerResponse := views.PlayerResponse{
		ID:      user.ID,
		Name:    user.FullName,
		Email:   user.Email,
		Ranking: user.Ranking,
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Player created successfully",
		Data:    playerResponse,
	})
}

// GetPlayer gets a User with player role by ID
func (pc *PlayerController) GetPlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid player ID",
		})
		return
	}

	var user models.User
	if err := pc.db.Where("role = ? AND id = ?", "player", uint(id)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Player not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch player",
		})
		return
	}

	playerResponse := views.PlayerResponse{
		ID:      user.ID,
		Name:    user.FullName,
		Email:   user.Email,
		Ranking: user.Ranking,
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Player retrieved successfully",
		Data:    playerResponse,
	})
}

// UpdatePlayer updates a User with player role
func (pc *PlayerController) UpdatePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid player ID",
		})
		return
	}

	var user models.User
	if err := pc.db.Where("role = ? AND id = ?", "player", uint(id)).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "not_found",
				Message: "Player not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch player",
		})
		return
	}

	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Ranking int    `json:"ranking"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Update user fields
	if req.Name != "" {
		user.FullName = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.Ranking = req.Ranking

	if err := pc.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update player",
		})
		return
	}

	playerResponse := views.PlayerResponse{
		ID:      user.ID,
		Name:    user.FullName,
		Email:   user.Email,
		Ranking: user.Ranking,
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Player updated successfully",
		Data:    playerResponse,
	})
}

// DeletePlayer soft deletes a User with player role
func (pc *PlayerController) DeletePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid player ID",
		})
		return
	}

	// Only allow deletion of users with player role
	result := pc.db.Where("role = ? AND id = ?", "player", uint(id)).Delete(&models.User{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to delete player",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, views.ErrorResponse{
			Error:   "not_found",
			Message: "Player not found",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Player deleted successfully",
	})
}
