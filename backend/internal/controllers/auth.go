package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"badminton-backend/internal/middleware"
	"badminton-backend/internal/models"
	"badminton-backend/internal/views"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

// Register creates new user account
func (ac *AuthController) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"full_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// All new users are players by default
	// Only admin can change user roles
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     models.RolePlayer, // Always player for registration
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "password_hash_failed",
			Message: "Failed to hash password",
		})
		return
	}

	// Save user
	if err := ac.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to create user account",
		})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusCreated, views.SuccessResponse{
		Message: "User registered successfully",
		Data: gin.H{
			"user":  views.ToUserResponse(user),
			"token": token,
		},
	})
}

// Login authenticates user
func (ac *AuthController) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Find user by username or email
	var user models.User
	err := ac.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, views.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid username or password",
		})
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, views.ErrorResponse{
			Error:   "invalid_credentials",
			Message: "Invalid username or password",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, views.ErrorResponse{
			Error:   "account_deactivated",
			Message: "Your account has been deactivated",
		})
		return
	}

	// Generate token
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "token_generation_failed",
			Message: "Failed to generate authentication token",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Login successful",
		Data: gin.H{
			"user":  views.ToUserResponse(user),
			"token": token,
		},
	})
}

// GetProfile returns current user profile
func (ac *AuthController) GetProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	c.JSON(http.StatusOK, views.ToUserResponse(*userObj))
}

// UpdateProfile updates current user profile
func (ac *AuthController) UpdateProfile(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Ranking  int    `json:"ranking"` // Only for players
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Update fields
	if req.FullName != "" {
		userObj.FullName = req.FullName
	}
	if req.Email != "" {
		userObj.Email = req.Email
	}
	if req.Ranking > 0 && userObj.IsPlayer() {
		userObj.Ranking = req.Ranking
	}

	if err := ac.db.Save(userObj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Profile updated successfully",
		Data:    views.ToUserResponse(*userObj),
	})
}

// ChangePassword changes user password
func (ac *AuthController) ChangePassword(c *gin.Context) {
	user, _ := c.Get("user")
	userObj := user.(*models.User)

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Verify current password
	if !userObj.CheckPassword(req.CurrentPassword) {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_current_password",
			Message: "Current password is incorrect",
		})
		return
	}

	// Update password
	userObj.Password = req.NewPassword
	if err := userObj.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "password_hash_failed",
			Message: "Failed to hash new password",
		})
		return
	}

	if err := ac.db.Save(userObj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Password changed successfully",
	})
}

// UpdateUserRole allows admin to change user role (Admin only)
func (ac *AuthController) UpdateUserRole(c *gin.Context) {
	userID := c.Param("user_id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_input",
			Message: err.Error(),
		})
		return
	}

	// Validate role
	if req.Role != string(models.RolePlayer) && req.Role != string(models.RoleAdmin) {
		c.JSON(http.StatusBadRequest, views.ErrorResponse{
			Error:   "invalid_role",
			Message: "Role must be 'player' or 'admin'",
		})
		return
	}

	// Find target user
	var targetUser models.User
	if err := ac.db.First(&targetUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, views.ErrorResponse{
			Error:   "user_not_found",
			Message: "User not found",
		})
		return
	}

	// Update role
	targetUser.Role = models.UserRole(req.Role)
	if err := ac.db.Save(&targetUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to update user role",
		})
		return
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "User role updated successfully",
		Data:    views.ToUserResponse(targetUser),
	})
}

// GetAllUsers returns all users (Admin only)
func (ac *AuthController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := ac.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, views.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to fetch users",
		})
		return
	}

	userResponses := make([]views.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = views.ToUserResponse(user)
	}

	c.JSON(http.StatusOK, views.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    userResponses,
	})
}
