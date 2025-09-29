package models

import "time"

// MatchStatus defines match statuses
type MatchStatus string

const (
	MatchPending   MatchStatus = "pending"
	MatchOngoing   MatchStatus = "ongoing"
	MatchCompleted MatchStatus = "completed"
	MatchCancelled MatchStatus = "cancelled"
)

// MatchType defines match types
type MatchType string

const (
	MatchSingles MatchType = "singles"
	MatchDoubles MatchType = "doubles"
)

// Match represents a badminton match
type Match struct {
	BaseModel
	TournamentID *uint       `json:"tournament_id"`
	Type         MatchType   `json:"type" gorm:"default:'singles'"`
	Status       MatchStatus `json:"status" gorm:"default:'pending'"`
	MatchDate    time.Time   `json:"match_date"`
	Round        string      `json:"round"` // qualification, round1, quarter, semi, final

	// For singles matches
	Player1ID    *uint `json:"player1_id"`
	Player2ID    *uint `json:"player2_id"`
	Player1Score int   `json:"player1_score" gorm:"default:0"`
	Player2Score int   `json:"player2_score" gorm:"default:0"`

	// For doubles matches
	Team1ID    *uint `json:"team1_id"`
	Team2ID    *uint `json:"team2_id"`
	Team1Score int   `json:"team1_score" gorm:"default:0"`
	Team2Score int   `json:"team2_score" gorm:"default:0"`

	// Winner tracking
	WinnerPlayerID *uint `json:"winner_player_id"`
	WinnerTeamID   *uint `json:"winner_team_id"`

	// Relations
	Tournament   *Tournament `json:"tournament,omitempty" gorm:"foreignKey:TournamentID"`
	Player1      *User       `json:"player1,omitempty" gorm:"foreignKey:Player1ID"`
	Player2      *User       `json:"player2,omitempty" gorm:"foreignKey:Player2ID"`
	Team1        *Team       `json:"team1,omitempty" gorm:"foreignKey:Team1ID"`
	Team2        *Team       `json:"team2,omitempty" gorm:"foreignKey:Team2ID"`
	WinnerPlayer *User       `json:"winner_player,omitempty" gorm:"foreignKey:WinnerPlayerID"`
	WinnerTeam   *Team       `json:"winner_team,omitempty" gorm:"foreignKey:WinnerTeamID"`
}

// IsTeamMatch checks if this is a team match
func (m *Match) IsTeamMatch() bool {
	return m.Type == MatchDoubles
}

// SetWinner sets the winner of the match
func (m *Match) SetWinner(playerID *uint, teamID *uint) {
	if m.IsTeamMatch() {
		m.WinnerTeamID = teamID
	} else {
		m.WinnerPlayerID = playerID
	}
	m.Status = MatchCompleted
}
