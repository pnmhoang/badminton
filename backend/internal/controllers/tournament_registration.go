package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"badminton-backend/internal/models"
	"badminton-backend/internal/views"
)

type TournamentRegistrationController struct {
	db *gorm.DB
}

func NewTournamentRegistrationController(db *gorm.DB) *TournamentRegistrationController {
	return &TournamentRegistrationController{db: db}
}

// RegisterForTournament allows players to register for tournaments
func (tc *TournamentRegistrationController) RegisterForTournament(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	if !userObj.IsPlayer() {
		c.JSON(http.StatusForbidden, views.ErrorResponse{
			Error:   "player_required",
			Message: "Only players can register for tournaments",
		})
		return
	}

	tournamentIDParam := c.Param("tournament_id")
	tournamentID, err := strconv.ParseUint(tournamentIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_tournament_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	// Check if tournament exists and is accepting registrations
	var tournament models.Tournament
	if err := tc.db.First(&tournament, uint(tournamentID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "tournament_not_found",
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

	if tournament.Status != models.TournamentUpcoming {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "registration_closed",
			Message: "Registration is closed for this tournament",
		})
		return
	}

	// Check if player is already registered
	var existingRegistration models.TournamentPlayer
	err = tc.db.Where("tournament_id = ? AND player_id = ?", tournamentID, userObj.ID).First(&existingRegistration).Error
	if err == nil {
		c.JSON(http.StatusConflict, views.ErrorResponse{
			Error:   "already_registered",
			Message: "You are already registered for this tournament",
		})
		return
	}

	// Check if tournament is full
	var registrationCount int64
	tc.db.Model(&models.TournamentPlayer{}).Where("tournament_id = ?", tournamentID).Count(&registrationCount)

	maxParticipants := tournament.GetMaxParticipants()
	if int(registrationCount) >= maxParticipants {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "tournament_full",
			Message: "Tournament is full",
		})
		return
	}

	// Create registration
	registration := models.TournamentPlayer{
		TournamentID: uint(tournamentID),
		PlayerID:     userObj.ID,
		Status:       "registered",
	}

	if err := tc.db.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to register for tournament",
		})
		return
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Successfully registered for tournament",
		Data:    gin.H{"registration_id": registration.ID},
	})
}

// UnregisterFromTournament allows players to withdraw from tournaments
func (tc *TournamentRegistrationController) UnregisterFromTournament(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	tournamentIDParam := c.Param("tournament_id")
	tournamentID, err := strconv.ParseUint(tournamentIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_tournament_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	// Find registration
	var registration models.TournamentPlayer
	err = tc.db.Where("tournament_id = ? AND player_id = ?", tournamentID, userObj.ID).First(&registration).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, views.ErrorResponse{
				Error:   "registration_not_found",
				Message: "You are not registered for this tournament",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to find registration",
		})
		return
	}

	// Update status to withdrawn instead of deleting
	registration.Status = "withdrawn"
	if err := tc.db.Save(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to withdraw from tournament",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Successfully withdrawn from tournament",
	})
}

// GetMyRegistrations returns current user's tournament registrations
func (tc *TournamentRegistrationController) GetMyRegistrations(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	var registrations []models.TournamentPlayer
	err := tc.db.Preload("Tournament").Where("player_id = ? AND status != ?", userObj.ID, "withdrawn").Find(&registrations).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch registrations",
		})
		return
	}

	c.JSON(http.StatusOK, registrations)
}

// CreateTeam creates a new team for doubles tournaments
func (tc *TournamentRegistrationController) CreateTeam(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		PartnerID   uint   `json:"partner_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Check if partner exists and is a player
	var partner models.User
	if err := tc.db.First(&partner, req.PartnerID).Error; err != nil {
		c.JSON(http.StatusNotFound, views.ErrorResponse{
			Error:   "partner_not_found",
			Message: "Partner not found",
		})
		return
	}

	if !partner.IsPlayer() {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "partner_not_player",
			Message: "Partner must be a player",
		})
		return
	}

	// Create team
	team := models.Team{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := tc.db.Create(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create team",
		})
		return
	}

	// Add team members
	teamPlayers := []models.TeamPlayer{
		{
			TeamID:   team.ID,
			PlayerID: userObj.ID,
			Role:     "captain",
		},
		{
			TeamID:   team.ID,
			PlayerID: partner.ID,
			Role:     "player",
		},
	}

	for _, tp := range teamPlayers {
		if err := tc.db.Create(&tp).Error; err != nil {
			c.JSON(http.StatusInternalServerError, views.ErrorResponse{
				Error:   "database_error",
				Message: "Failed to add team members",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Team created successfully",
		Data:    gin.H{"team_id": team.ID, "team_name": team.Name},
	})
}

// RegisterTeamForTournament registers a team for doubles tournament
func (tc *TournamentRegistrationController) RegisterTeamForTournament(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	tournamentIDParam := c.Param("tournament_id")
	tournamentID, err := strconv.ParseUint(tournamentIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_tournament_id",
			Message: "Invalid tournament ID",
		})
		return
	}

	var req struct {
		TeamID uint `json:"team_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Check if tournament is doubles type
	var tournament models.Tournament
	if err := tc.db.First(&tournament, uint(tournamentID)).Error; err != nil {
		c.JSON(http.StatusNotFound, views.ErrorResponse{
			Error:   "tournament_not_found",
			Message: "Tournament not found",
		})
		return
	}

	if !tournament.IsTeamTournament() {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "singles_tournament",
			Message: "This is a singles tournament, not doubles",
		})
		return
	}

	// Check if user is member of the team
	var teamPlayer models.TeamPlayer
	if err := tc.db.Where("team_id = ? AND player_id = ?", req.TeamID, userObj.ID).First(&teamPlayer).Error; err != nil {
		c.JSON(http.StatusForbidden, views.ErrorResponse{
			Error:   "not_team_member",
			Message: "You are not a member of this team",
		})
		return
	}

	// Check if team is already registered
	var existingRegistration models.TournamentTeam
	err = tc.db.Where("tournament_id = ? AND team_id = ?", tournamentID, req.TeamID).First(&existingRegistration).Error
	if err == nil {
		c.JSON(http.StatusConflict, views.ErrorResponse{
			Error:   "team_already_registered",
			Message: "Team is already registered for this tournament",
		})
		return
	}

	// Create team registration
	registration := models.TournamentTeam{
		TournamentID: uint(tournamentID),
		TeamID:       req.TeamID,
		Status:       "registered",
	}

	if err := tc.db.Create(&registration).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to register team for tournament",
		})
		return
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "Team successfully registered for tournament",
		Data:    gin.H{"registration_id": registration.ID},
	})
}
