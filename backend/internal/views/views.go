package views

import "badminton-backend/internal/models"

type PlayerResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Ranking int    `json:"ranking"`
}

type MatchResponse struct {
	ID           uint            `json:"id"`
	Player1      PlayerResponse  `json:"player1"`
	Player2      PlayerResponse  `json:"player2"`
	Player1Score int             `json:"player1_score"`
	Player2Score int             `json:"player2_score"`
	Status       string          `json:"status"`
	MatchDate    string          `json:"match_date"`
	Tournament   *TournamentInfo `json:"tournament,omitempty"`
}

type TournamentResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
	MaxPlayers  int    `json:"max_players"`
	MatchCount  int    `json:"match_count"`
}

type TournamentInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
	Ranking  int    `json:"ranking,omitempty"` // Only for players
}

// Helper functions to convert models to responses
func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     string(user.Role),
		IsActive: user.IsActive,
		Ranking:  user.Ranking,
	}
}

func ToPlayerResponse(player models.Player) PlayerResponse {
	return PlayerResponse{
		ID:      player.ID,
		Name:    player.Name,
		Email:   player.Email,
		Ranking: player.Ranking,
	}
}

func ToMatchResponse(match models.Match) MatchResponse {
	var player1, player2 PlayerResponse
	if match.Player1 != nil {
		player1 = PlayerResponse{
			ID:      match.Player1.ID,
			Name:    match.Player1.FullName,
			Email:   match.Player1.Email,
			Ranking: match.Player1.Ranking,
		}
	}
	if match.Player2 != nil {
		player2 = PlayerResponse{
			ID:      match.Player2.ID,
			Name:    match.Player2.FullName,
			Email:   match.Player2.Email,
			Ranking: match.Player2.Ranking,
		}
	}

	response := MatchResponse{
		ID:           match.ID,
		Player1:      player1,
		Player2:      player2,
		Player1Score: match.Player1Score,
		Player2Score: match.Player2Score,
		Status:       string(match.Status),
		MatchDate:    match.MatchDate.Format("2006-01-02 15:04:05"),
	}

	if match.Tournament != nil {
		response.Tournament = &TournamentInfo{
			ID:   match.Tournament.ID,
			Name: match.Tournament.Name,
		}
	}

	return response
}

func ToTournamentResponse(tournament models.Tournament) TournamentResponse {
	return TournamentResponse{
		ID:          tournament.ID,
		Name:        tournament.Name,
		Description: tournament.Description,
		StartDate:   tournament.StartDate.Format("2006-01-02"),
		EndDate:     tournament.EndDate.Format("2006-01-02"),
		Status:      string(tournament.Status),
		MaxPlayers:  tournament.GetMaxParticipants(),
		MatchCount:  len(tournament.Matches),
	}
}
