package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"badminton-backend/internal/controllers"
	"badminton-backend/internal/middleware"
	"badminton-backend/internal/models"
)

func main() {
	// Initialize database
	db, err := gorm.Open(sqlite.Open("badminton.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	db.AutoMigrate(
		&models.User{},
		&models.Player{}, // Keep for backward compatibility
		&models.Team{},
		&models.TeamPlayer{},
		&models.Tournament{},
		&models.TournamentPlayer{},
		&models.TournamentTeam{},
		&models.Match{},
	)

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Initialize controllers
	authController := controllers.NewAuthController(db)
	playerController := controllers.NewPlayerController(db)
	matchController := controllers.NewMatchController(db)
	tournamentController := controllers.NewTournamentController(db)
	tournamentRegController := controllers.NewTournamentRegistrationController(db)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		v1.POST("/register", authController.Register)
		v1.POST("/login", authController.Login)

		// Protected routes
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware(db))
		{
			// Profile routes
			authorized.GET("/profile", authController.GetProfile)
			authorized.PUT("/profile", authController.UpdateProfile)
			authorized.POST("/change-password", authController.ChangePassword)

			// Admin-only user management routes
			authorized.GET("/users", middleware.RequireAdmin(), authController.GetAllUsers)
			authorized.PUT("/users/:user_id/role", middleware.RequireAdmin(), authController.UpdateUserRole)

			// Player routes (keep for backward compatibility)
			authorized.GET("/players", playerController.GetPlayers)
			authorized.POST("/players", playerController.CreatePlayer)
			authorized.GET("/players/:id", playerController.GetPlayer)
			authorized.PUT("/players/:id", playerController.UpdatePlayer)
			authorized.DELETE("/players/:id", playerController.DeletePlayer)

			// Match routes
			authorized.GET("/matches", matchController.GetMatches)
			authorized.POST("/matches", matchController.CreateMatch)
			authorized.GET("/matches/:id", matchController.GetMatch)
			authorized.PUT("/matches/:id", matchController.UpdateMatch)
			authorized.DELETE("/matches/:id", matchController.DeleteMatch)

			// Tournament routes
			authorized.GET("/tournaments", tournamentController.GetTournaments)
			authorized.POST("/tournaments", middleware.RequireAdmin(), tournamentController.CreateTournament)
			authorized.GET("/tournaments/:id", tournamentController.GetTournament)
			authorized.PUT("/tournaments/:id", middleware.RequireAdmin(), tournamentController.UpdateTournament)
			authorized.DELETE("/tournaments/:id", middleware.RequireAdmin(), tournamentController.DeleteTournament)

			// Tournament registration routes
			authorized.POST("/tournament-registration/:tournament_id", middleware.RequirePlayer(), tournamentRegController.RegisterForTournament)
			authorized.DELETE("/tournament-registration/:tournament_id", middleware.RequirePlayer(), tournamentRegController.UnregisterFromTournament)
			authorized.GET("/my-registrations", middleware.RequirePlayer(), tournamentRegController.GetMyRegistrations)

			// Team routes
			authorized.POST("/teams", middleware.RequirePlayer(), tournamentRegController.CreateTeam)
			authorized.POST("/team-registration/:tournament_id", middleware.RequirePlayer(), tournamentRegController.RegisterTeamForTournament)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start server
	log.Println("Server starting on :8080")
	r.Run(":8080")
}
